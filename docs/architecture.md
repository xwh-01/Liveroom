# 架构设计文档

## 整体架构

```
┌─────────────────────────────────────────────────┐
│                  Vue Frontend                    │
│  (LivePlayer / DanmuPanel / GiftPanel /         │
│   RankBoard / RoomStats / SystemLog)             │
└──────────┬───────────────┬──────────────────────┘
           │ WebSocket     │ HTTP
           ▼               ▼
┌──────────────────────────────────────────────────┐
│                  Gin Router                       │
│  /ws  /api/room/state  /api/room/rank            │
│  /api/room/chats  /api/room/gifts                │
└──────────┬───────────────┬──────────────────────┘
           │               │
           ▼               ▼
┌──────────────────┐  ┌──────────────────┐
│  WSController    │  │  RoomController  │
│  (readPump/      │  │  (HTTP handlers) │
│   writePump)     │  │                  │
└────────┬─────────┘  └────────┬─────────┘
         │                     │
         ▼                     ▼
┌──────────────────────────────────────────────────┐
│              MessageDispatcher                    │
│  chat -> ChatService.HandleChat                  │
│  gift -> GiftService.HandleGift                   │
└────────┬─────────────────────────────────────────┘
         │
    ┌────┴─────────────────────┐
    ▼                          ▼
┌───────────┐            ┌───────────┐
│ChatService│            │GiftService│
│           │            │           │
│ 1.限流检查│            │ 1.礼物校验│
│ 2.广播弹幕│            │ 2.Redis排行│
│ 3.IncrCnt │            │ 3.广播礼物│
│ 4.Submit  │            │ 4.广播排行│
│   persist │            │ 5.IncrCnt │
└─────┬─────┘            │ 6.Submit  │
      │                  │   persist │
      │                  └─────┬─────┘
      │                        │
      │  Submit (non-blocking) │
      ▼                        ▼
┌─────────────────────────────────────┐
│         PersistService              │
│  ┌─────────────────────────────┐   │
│  │  queue (chan, cap=10000)    │   │
│  │  ┌─────┬─────┬─────┬─────┐ │   │
│  │  │Event│Event│Event│ ... │ │   │
│  │  └──┬──┴──┬──┴──┬──┴─────┘ │   │
│  │     │     │     │           │   │
│  │  ┌──▼─────▼─────▼──┐       │   │
│  │  │  Worker (x2)     │       │   │
│  │  │  handle(chat) -> RecordDao.InsertChatRecord    │
│  │  │  handle(gift) -> RecordDao.InsertGiftRecord    │
│  │  └────────┬────────┘       │   │
│  └───────────┼────────────────┘   │
│              │                    │
│  droppedCount (atomic.Int64)     │
└──────────────┼────────────────────┘
               │
               ▼
┌──────────────────────┐
│     RecordDao        │
│  database/sql        │
│  INSERT / SELECT     │
└──────────┬───────────┘
           │
           ▼
┌──────────────────────┐
│       MySQL          │
│  chat_records        │
│  gift_records        │
└──────────────────────┘
```

## 数据流

### 弹幕消息流

```
Client -> WebSocket readPump -> MessageDispatcher.Dispatch("chat")
  -> ChatService.HandleChat
    -> RateLimitService.IsLimited (Redis INCR + EXPIRE)
    -> [限流] SendToUser system message (仅通知自己)
    -> RoomHub.Broadcast chat -> 所有客户端 writePump -> Client
    -> RedisDao.IncrChatCount (Redis INCR)
    -> PersistService.Submit(chat event) -> queue -> worker -> RecordDao.InsertChatRecord -> MySQL chat_records
```

### 礼物流水

```
Client -> WebSocket readPump -> MessageDispatcher.Dispatch("gift")
  -> GiftService.HandleGift
    -> GetGiftScore (校验礼物类型)
    -> RedisDao.AddGiftScore (Redis ZINCRBY)
    -> RoomHub.Broadcast gift -> 所有客户端 writePump -> Client
    -> RedisDao.GetTopRank (Redis ZREVRANGE)
    -> RoomHub.Broadcast rank -> 所有客户端 writePump -> Client
    -> RedisDao.IncrGiftCount (Redis INCR)
    -> PersistService.Submit(gift event) -> queue -> worker -> RecordDao.InsertGiftRecord -> MySQL gift_records
```

### 异步落库流程

```
ChatService/GiftService
  -> PersistService.Submit(event)
    -> 非阻塞写入 chan (select + default)
    -> 队列满: droppedCount++ + slog.Warn + return false
  -> worker goroutine
    -> 从 chan 消费
    -> 按 Type 分发:
       "chat" -> RecordDao.InsertChatRecord
       "gift" -> RecordDao.InsertGiftRecord
    -> 写入失败: slog.Error + 继续循环
```

## 依赖关系

```
main.go (wire)
 ├── config.Load
 ├── common.InitRedis -> RedisDao
 ├── common.InitMySQL -> RecordDao
 ├── hub.NewRoomHub
 ├── service.NewPersistService(recordDao, 10000)
 │    └── persistSvc.Start(ctx, 2)
 ├── service.NewRoomService(redisDao, roomHub, persistSvc)
 ├── service.NewChatService(rateLimitSvc, roomHub, redisDao, persistSvc)
 ├── service.NewGiftService(redisDao, roomHub, persistSvc)
 ├── controller.NewWSController(dispatcher, roomSvc)
 └── controller.NewRoomController(roomSvc, rankSvc, recordDao)
```

## MySQL 表结构

### chat_records

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| room_id | VARCHAR(64) | 房间 ID |
| user_id | VARCHAR(64) | 用户 ID |
| content | VARCHAR(512) | 弹幕内容 |
| created_at | DATETIME | 创建时间 |

索引：`idx_room_created_at (room_id, created_at)`, `idx_user_created_at (user_id, created_at)`

### gift_records

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键，自增 |
| room_id | VARCHAR(64) | 房间 ID |
| user_id | VARCHAR(64) | 用户 ID |
| gift_type | VARCHAR(32) | 礼物类型 (heart/rocket) |
| gift_score | INT | 礼物分值 |
| created_at | DATETIME | 创建时间 |

索引：`idx_room_created_at (room_id, created_at)`, `idx_user_created_at (user_id, created_at)`

## 队列满处理

- PersistService.queue 大小为 10000（可配置）
- Submit 使用 `select { case ch <- event: return true; default: ... }` 非阻塞模式
- 队满时 `droppedCount` 自增（atomic.Int64）
- 打印 `slog.Warn` 日志，包含 event 类型和当前 dropped_total
- 通过 `/api/room/state` 的 `persist_dropped_count` 字段暴露
