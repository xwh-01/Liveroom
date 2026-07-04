# LiveRoom Battle

直播间实时互动系统 —— 基于 WebSocket + Redis + MySQL 实现弹幕、礼物、排行榜、限流、在线人数统计及历史查询。

**注意：本项目不包含真实音视频直播功能（无 RTMP / HLS / WebRTC），页面中仅为假直播画面。**

## 技术栈

**后端**
- Go 1.21+
- Gin (HTTP 路由)
- gorilla/websocket (WebSocket)
- Redis (限流、排行榜、在线人数)
- MySQL + GORM (弹幕记录、礼物流水、房间信息持久化)
- TOML 配置
- slog 日志

**前端**
- Vue 3
- Vue Router
- Element Plus
- Axios
- 原生 WebSocket API

## 目录结构

```
LiveRoom-Battle/
├── backend/
│   ├── main.go                  # 入口，依赖注入
│   ├── go.mod / go.sum
│   ├── common/                  # 基础设施初始化
│   │   ├── redis.go             # Redis 连接
│   │   ├── mysql.go             # MySQL 连接 + GORM 连接池
│   │   └── websocket.go         # WebSocket upgrader
│   ├── config/
│   │   ├── config.go            # 配置结构体
│   │   └── config.toml          # TOML 配置文件
│   ├── controller/              # HTTP / WebSocket 入口
│   │   ├── ws_controller.go     # WebSocket 处理器
│   │   └── room_controller.go   # 房间状态 / 排行榜 / 历史查询 HTTP 接口
│   ├── service/                 # 业务逻辑层
│   │   ├── chat_service.go      # 弹幕处理 (含 MySQL 落库)
│   │   ├── gift_service.go      # 礼物处理 (含 MySQL 落库)
│   │   ├── rank_service.go      # 排行榜
│   │   ├── room_service.go      # 房间管理 (含自动创建房间)
│   │   ├── rate_limit_service.go # 限流
│   │   ├── dispatcher.go        # 消息分发器
│   │   └── hub_interface.go     # Hub 接口定义
│   ├── dao/                     # 数据访问层
│   │   ├── redis_dao.go         # Redis 操作封装
│   │   ├── room_dao.go          # 房间 CRUD
│   │   ├── chat_message_dao.go  # 弹幕持久化
│   │   └── gift_record_dao.go   # 礼物流水持久化
│   ├── model/                   # 数据结构
│   │   ├── message.go           # WS 消息类型
│   │   ├── room.go              # 房间状态 (API 响应)
│   │   ├── room_entity.go       # 房间实体 (MySQL)
│   │   ├── chat_message.go      # 弹幕记录 (MySQL)
│   │   ├── gift_record.go       # 礼物记录 (MySQL)
│   │   ├── gift.go              # 礼物配置
│   │   └── client.go            # WebSocket 客户端
│   ├── hub/
│   │   └── room_hub.go          # 单机内存房间管理器
│   ├── router/
│   │   └── router.go            # 路由注册
│   ├── middleware/
│   │   └── cors.go              # CORS 中间件
│   ├── utils/
│   │   ├── response.go          # 统一响应
│   │   ├── keygen.go            # Redis key 生成
│   │   ├── time.go              # 时间工具
│   │   └── uuid.go              # UUID 生成
│   └── bot/
│       ├── go.mod
│       └── main.go              # Bot 模拟用户脚本
├── vue-frontend/                # Vue 3 前端
│   ├── src/
│   │   ├── views/LiveRoom.vue   # 主页面
│   │   ├── components/          # 组件
│   │   │   ├── LivePlayer.vue   # 假直播画面
│   │   │   ├── DanmuPanel.vue   # 弹幕区
│   │   │   ├── GiftPanel.vue    # 礼物按钮
│   │   │   ├── RankBoard.vue    # 排行榜
│   │   │   ├── RoomStats.vue    # 房间状态
│   │   │   └── SystemLog.vue    # 系统消息
│   │   ├── utils/ws.js          # WebSocket 客户端
│   │   └── router/index.js      # 路由
│   ├── package.json
│   └── vite.config.js
└── README.md
```

## 启动方式

### 前提条件
- Go 1.21+
- Redis（监听 127.0.0.1:6379）
- MySQL（监听 127.0.0.1:3306）
- Node.js 18+

### 1. 启动基础设施

确保 Redis 和 MySQL 已启动。MySQL 需要先创建数据库：

```sql
CREATE DATABASE IF NOT EXISTS liveroom DEFAULT CHARSET utf8mb4;
```

可修改 `backend/config/config.toml` 中的连接地址和账密。

### 2. 启动后端

```bash
cd backend
go run main.go
```

服务启动在 `http://localhost:8080`，启动时会自动迁移 MySQL 表结构（AutoMigrate）。

### 3. 启动前端

```bash
cd vue-frontend
npm install
npm run dev
```

前端启动在 `http://localhost:3000`，通过 Vite 代理转发 API 到后端。

### 4. 运行 Bot 模拟用户

```bash
cd backend/bot
go run main.go
```

Bot 会模拟 20 个用户连接直播间，随机发送弹幕和礼物。其中 4 个高频用户会触发限流。

## V2 MySQL 持久化

### 数据表

V2 新增 3 张 MySQL 表，启动时通过 GORM AutoMigrate 自动创建：

| 表 | 说明 |
|---|------|
| `rooms` | 直播间信息 (room_id, title, anchor_name, status) |
| `chat_messages` | 弹幕记录 (message_id, room_id, user_id, content) |
| `gift_records` | 礼物流水 (record_id, room_id, user_id, gift_id, score) |

### 弹幕保存链路

```
用户发送弹幕 → WebSocket → RateLimit 检查 → Redis 限流判断
→ [未超限] → 广播到房间 → 同步写入 MySQL chat_messages 表
→ [超限]   → 仅通知本人，不入库
```

### 礼物流水保存链路

```
用户送礼 → WebSocket → 验证礼物类型 → Redis ZIncrBy 更新积分
→ 广播礼物消息 → 查询 Redis 排行榜 Top10 → 广播排行
→ 同步写入 MySQL gift_records 表
```

### 房间自动创建

用户首次加入直播间时，检测 `room_id` 是否存在，不存在则自动创建默认房间记录。

### 历史查询接口

| 接口 | 说明 |
|------|------|
| `GET /api/v1/rooms/:room_id/messages?limit=50` | 查询最近弹幕历史，按时间倒序 |
| `GET /api/v1/rooms/:room_id/gifts?limit=50` | 查询最近礼物流水，按时间倒序 |

前端页面底部的「历史弹幕」和「礼物流水」按钮可调用上述接口。

> **注意**：当前 V2 使用同步落库。V3 将改用 RabbitMQ 异步落库，解耦写入路径。

## MVP 功能

### 1. WebSocket 直播间连接
- GET `/ws?room_id=xxx&user_id=xxx`
- 同一 room_id 下的用户属于同一房间
- 进入/离开时广播在线人数变化
- 读写分离，避免 concurrent write

### 2. 弹幕功能
- 发送 `chat` 类型消息
- 基于 Redis 的限流：同一用户同一房间 1 秒最多 5 条
- 触发限流只通知当前用户，不广播
- 限流时 `limited_count` 加 1

### 3. 礼物功能
- 支持 `heart`（10分）和 `rocket`（100分）
- Redis ZSet 维护榜单
- 送礼后广播礼物消息 + 最新 TOP10 排行榜

### 4. 排行榜
- HTTP GET `/api/room/rank?room_id=xxx`
- Redis ZSet 实现，按积分倒序

### 5. 房间状态
- HTTP GET `/api/room/state?room_id=xxx`
- 返回 online_count、limited_count

## WebSocket 消息协议

客户端发送：
```json
{ "type": "chat", "room_id": "1001", "user_id": "user1", "data": { "content": "666" } }
{ "type": "gift", "room_id": "1001", "user_id": "user1", "data": { "gift_type": "heart" } }
```

服务端推送：
```json
{ "type": "chat", "room_id": "1001", "user_id": "user1", "data": { "content": "666", "timestamp": "23:30:00" } }
{ "type": "gift", "room_id": "1001", "user_id": "user1", "data": { "gift_type": "heart", "gift_score": 10, "sender": "user1" } }
{ "type": "rank", "room_id": "1001", "data": { "rankings": [...] } }
{ "type": "online", "room_id": "1001", "data": { "count": 5 } }
{ "type": "system", "data": { "content": "你发送太快了，已被限流" } }
```

## 架构设计要点

- **分层清晰**：controller → service → dao，各层职责明确
- **依赖注入**：服务之间通过构造函数注入，不使用全局变量
- **消息分发器**：可扩展的消息类型注册机制，方便增加新消息类型
- **Redis key 统一管理**：通过 `utils/keygen.go` 集中生成
- **礼物配置集中**：礼物类型和分值在 `model/gift.go` 中管理
- **RoomHub 接口化**：内存房间管理器有清晰接口，后续可升级为分布式
- **WebSocket 读写分离**：每个连接独立的 readPump / writePump 协程

## 后续扩展路线

| 版本 | 功能 | 状态 |
|------|------|------|
| V1 | WebSocket + Redis 排行榜 + Redis 限流 | Done |
| V2 | MySQL 保存弹幕记录和礼物流水 | Done |
| V3 | RabbitMQ 异步落库，解耦写入 | TODO |
| V4 | Prometheus + Grafana 监控在线人数、QPS、限流次数 | TODO |
| V5 | 微服务拆分：chat-service、gift-service、rank-service | TODO |
| V6 | 多实例 WebSocket，Redis Pub/Sub 跨节点广播 | TODO |

## 当前不包含

- RabbitMQ / Kafka 消息队列
- 微服务拆分 / 服务注册发现
- 登录注册 / 权限系统
- 真实音视频直播（RTMP / HLS / WebRTC）
- AI Agent

## License

MIT
