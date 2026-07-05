# 架构设计文档

## 整体架构

```
┌─────────────────────────────────────────────────┐
│                  Vue Frontend                    │
│  (RoomLobby / LiveRoom / LivePlayer /           │
│   DanmuPanel / GiftPanel / RankBoard /          │
│   RoomStats / SystemLog)                        │
└──────────┬───────────────┬──────────────────────┘
           │ WebSocket     │ HTTP
           ▼               ▼
┌──────────────────────────────────────────────────┐
│                  Gin Router                       │
│  /ws  /api/room/state  /api/room/rank            │
│  /api/rooms  /api/rooms/:room_id                 │
│  /api/admin/rooms/:room_id/close                 │
└──────────┬───────────────┬──────────────────────┘
           │               │
           ▼               ▼
┌──────────────────┐  ┌──────────────────┐
│  WSController    │  │  RoomController  │
│  (readPump/      │  │  + RoomManage    │
│   writePump)     │  │    Service       │
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
└───────────┘            │ 4.广播排行│
                         │ 5.IncrCnt │
                         └───────────┘
                              │
                              ▼
                    ┌──────────────────┐
                    │  RoomHub (内存)   │
                    │  Broadcast        │
                    │  OnlineCount      │
                    └──────────────────┘
```

## 数据流

### 房间大厅流程

```
Browser -> GET /api/rooms
  -> RoomController.ListRooms
    -> RoomManageService.ListLiveRooms
      -> RedisDao.ListLiveRooms (ZREVRANGE room:live)
      -> 对每个 room_id:
        -> RedisDao.GetRoomMeta (HGETALL room:meta:{id})
        -> Hub.OnlineCount (内存)
        -> RedisDao.GetChatCount / GetGiftCount
      -> 返回排序后的房间列表
  -> 前端 RoomLobby 展示卡片
```

### 进入房间流程

```
Browser -> GET /api/rooms/:room_id
  -> RoomController.GetRoom
    -> RoomManageService.GetRoom
      -> RedisDao.GetRoomMeta (HGETALL room:meta:{id})
        -> 不存在 -> 404
        -> status=closed -> 前端提示"房间已关闭"
        -> status=live -> 返回房间详情
  -> 前端 LiveRoom 获取信息后:
    -> WebSocket connect /ws?room_id=xxx&user_id=xxx
      -> WSController 检查 room_meta (存在 + status=live)
      -> RoomService.Join -> RoomHub + Redis online set
      -> 广播 online 消息
```

### 弹幕消息流

```
Client -> WebSocket readPump -> MessageDispatcher.Dispatch("chat")
  -> ChatService.HandleChat
    -> RateLimitService.IsLimited (Redis INCR + EXPIRE)
    -> [限流] SendToUser system message (仅通知自己)
    -> RoomHub.Broadcast chat -> 所有客户端 writePump -> Client
    -> RedisDao.IncrChatCount (Redis INCR)
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
```

## 依赖关系

```
main.go (wire)
 ├── config.Load
 ├── common.InitRedis -> RedisDao
 ├── hub.NewRoomHub
 ├── service.NewRoomService(redisDao, roomHub)
 ├── service.NewRoomManageService(redisDao, roomHub)
 │    └── roomManageSvc.EnsureDefaultRooms(ctx)
 ├── service.NewChatService(rateLimitSvc, roomHub, redisDao)
 ├── service.NewGiftService(redisDao, roomHub)
 ├── controller.NewWSController(dispatcher, roomSvc, roomManageSvc)
 └── controller.NewRoomController(roomSvc, roomManageSvc, rankSvc)
```

## 预置房间

后端启动时 `EnsureDefaultRooms` 通过 Redis Hash 写入 3 个默认房间：

| room_id | title | anchor_name |
|---------|-------|-------------|
| 1001 | 游戏开黑直播间 | 小黑 |
| 1002 | 音乐闲聊直播间 | Echo |
| 1003 | 学习摸鱼直播间 | 阿码 |

如果房间已存在则跳过，保证 bot 默认 room_id=1001 可用。
