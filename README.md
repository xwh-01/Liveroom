# LiveRoom Battle

直播间实时互动系统 V3 —— 基于 Go + WebSocket + Redis + MySQL 实现预置房间大厅、动态进入房间、弹幕、礼物、排行榜、限流、在线人数统计、bot 压测和异步持久化。

当前能力：
- 预置房间大厅、动态进入房间
- WebSocket 弹幕、礼物广播
- Redis 礼物排行榜、限流、在线人数
- bot 压测
- MySQL 弹幕记录、礼物流水
- PersistService 本机异步落库队列
- /api/room/persist/state 查看落库状态

当前不包含：RabbitMQ、微服务、登录注册、主播后台、用户创建房间、真实音视频直播。

**注意：本项目不包含真实音视频直播功能（无 RTMP / HLS / WebRTC），页面中仅为假直播画面。**

## 技术栈

**后端**
- Go 1.21+
- Gin (HTTP 路由)
- gorilla/websocket (WebSocket)
- Redis (限流、排行榜、在线人数、房间元数据)
- MySQL (弹幕记录、礼物流水持久化)
- database/sql + go-sql-driver/mysql
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
│   │   └── websocket.go         # WebSocket upgrader
│   ├── config/
│   │   ├── config.go            # 配置结构体
│   │   └── config.toml          # TOML 配置文件
│   ├── controller/              # HTTP / WebSocket 入口
│   │   ├── ws_controller.go     # WebSocket 处理器
│   │   └── room_controller.go   # HTTP 接口
│   ├── service/                 # 业务逻辑层
│   │   ├── chat_service.go      # 弹幕处理
│   │   ├── gift_service.go      # 礼物处理
│   │   ├── rank_service.go      # 排行榜
│   │   ├── room_service.go      # 房间管理
│   │   ├── room_manage_service.go # 房间大厅
│   │   ├── rate_limit_service.go # 限流
│   │   └── dispatcher.go        # 消息分发器
│   ├── dao/                     # 数据访问层
│   │   └── redis_dao.go         # Redis 操作封装
│   ├── model/                   # 数据结构
│   │   ├── message.go           # WS 消息类型
│   │   ├── room.go              # 房间状态
│   │   ├── room_meta.go         # 房间元数据
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
│   │   └── time.go              # 时间工具
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
├── docs/
│   ├── benchmark.md             # 压测指南
│   ├── benchmark-result.md      # 压测结果模板
│   ├── architecture.md          # 架构设计文档
│   └── demo.md                  # 功能演示验证指南
└── README.md
```

## 启动方式

### 前提条件
- Go 1.21+
- Redis（监听 127.0.0.1:6379）
- MySQL（监听 127.0.0.1:3306）
- Node.js 18+

### 1. 初始化 MySQL

```bash
mysql -uroot -p < backend/sql/001_init_records.sql
```

可修改 `backend/config/config.toml` 中的 MySQL 连接地址和账号密码。

### 2. 启动 Redis

```bash
redis-server
```

### 3. 启动后端

```bash
cd backend
go run main.go
```

服务启动在 `http://localhost:8080`。

### 4. 启动前端

```bash
cd vue-frontend
npm install
npm run dev
```

前端启动在 `http://localhost:3000`，通过 Vite 代理转发 API 到后端。

### 5. 运行 Bot 模拟用户

```bash
cd backend/bot
go run main.go
```

默认参数：20 用户、60 秒、chat 间隔 500ms、gift 间隔 1500ms。

```bash
# 自定义参数
go run main.go -host=localhost:8080 -room_id=1001 -user_count=50 -duration_seconds=120 -chat_interval_ms=200 -gift_interval_ms=800
```

详细参数见 `docs/benchmark.md`。

## 功能

### 1. WebSocket 直播间连接
- GET `/ws?room_id=xxx&user_id=xxx`
- 同一 room_id 下的用户属于同一房间
- 进入/离开时广播在线人数变化
- 读写分离，避免 concurrent write

### 2. 弹幕功能
- 发送 `chat` 类型消息
- 基于 Redis 的限流：同一用户同一房间 1 秒最多 5 条
- 触发限流只通知当前用户，不广播
- 限流时 `limited_count` 加 1（Redis 持久化）

### 3. 礼物功能
- 支持 `heart`（10分）和 `rocket`（100分）
- Redis ZSet 维护榜单
- 送礼后广播礼物消息 + 最新 TOP10 排行榜

### 4. 排行榜
- HTTP GET `/api/room/rank?room_id=xxx`
- Redis ZSet 实现，按积分倒序

### 5. 房间状态
- HTTP GET `/api/room/state?room_id=xxx`
- 返回字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `room_id` | string | 房间 ID |
| `online_count` | int | 当前在线连接数（来自 Hub 内存） |
| `limited_count` | int | 累计限流次数（来自 Redis） |
| `chat_count` | int64 | 累计弹幕数（来自 Redis） |
| `gift_count` | int64 | 累计礼物数（来自 Redis） |

V2 房间大厅 API：

- `GET /api/rooms?limit=20` — 获取直播间列表
- `GET /api/rooms/:room_id` — 获取单个房间信息
- `POST /api/admin/rooms/:room_id/close` — 关闭房间（开发环境用，无鉴权）

V3 历史记录 API：

- `GET /api/room/chats?room_id=xxx&limit=20` — 查询最近弹幕
- `GET /api/room/gifts?room_id=xxx&limit=20` — 查询最近礼物流水
- `GET /api/room/persist/state` — 查看异步落库队列状态

### 6. 可观测指标

| 指标 | 存储 | 暴露方式 |
|------|------|----------|
| 在线人数 | Hub 内存 | `/api/room/state` + WS `online` 推送 |
| 限流次数 | Redis | `/api/room/state` |
| chat_count | Redis | `/api/room/state` |
| gift_count | Redis | `/api/room/state` |
| 广播延迟 | slog 日志 | >10ms 时打印 `latency_us` |
| 消息丢弃 | slog 日志 | send buffer 满时 Warn |

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
- **广播可观测**：bot 模拟多用户并发接入，广播日志记录 target_clients、dropped_clients、latency_us

## 后续扩展路线

| 版本 | 功能 | 状态 |
|------|------|------|
| V1 | WebSocket + Redis 排行榜 + Redis 限流 | Done |
| V1.1 | 内存可观测指标、增强 /api/room/state、bot CLI 参数、压测指南 | Done |
| V1.2 | Broadcast 可观测日志、压测结果模板 benchmark-result.md | Done |
| V2 | 预置房间大厅 + 动态进入房间 | Done |
| V3 | MySQL 异步持久化弹幕记录和礼物流水 | Done |
| V4 | RabbitMQ 异步落库 | TODO |
| V5 | Prometheus + Grafana 监控 | TODO |
| V6 | 多实例 WebSocket，Redis Pub/Sub 跨节点广播 | TODO |

## 当前不包含

- 普通用户创建房间
- 主播后台
- 登录注册 / 权限系统
- 真实音视频直播（RTMP / HLS / WebRTC）
- RabbitMQ / Kafka 消息队列
- 微服务拆分 / 服务注册发现
- AI Agent

## 文档

- [docs/benchmark.md](docs/benchmark.md) — 压测指南
- [docs/benchmark-result.md](docs/benchmark-result.md) — 压测结果模板
- [docs/architecture.md](docs/architecture.md) — 架构设计文档
- [docs/demo.md](docs/demo.md) — 功能演示验证指南

## License

MIT
