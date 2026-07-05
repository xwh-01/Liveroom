# LiveRoom Battle 压测指南

## 环境准备

### 1. 启动 Redis

```bash
redis-server
```

默认监听 `localhost:6379`，无需密码。

### 2. 启动后端

```bash
cd backend
go run main.go
```

服务启动在 `http://localhost:8080`。

### 3. 启动前端（可选）

```bash
cd vue-frontend
npm install
npm run dev
```

页面打开 `http://localhost:3000`，输入 Room ID 和 User ID 后点连接。

### 4. 启动 bot 压测

```bash
cd backend/bot
go run main.go
```

默认参数：20 用户连接 `room 1001`，运行 60 秒。

#### Bot 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-host` | `localhost:8080` | 服务器地址 |
| `-room_id` | `1001` | 目标房间 ID |
| `-user_count` | `20` | 模拟用户数 |
| `-duration_seconds` | `60` | 压测持续时间（秒） |
| `-chat_interval_ms` | `500` | 每个用户发弹幕间隔（毫秒） |
| `-gift_interval_ms` | `1500` | 每个用户送礼间隔（毫秒） |

Bot 使用独立的 chat ticker 和 gift ticker，按各自间隔独立发送，支持精确控制压测参数。

运行结束后输出：连接成功人数、连接失败人数、chat 发送数、gift 发送数、写入错误数、断开连接数、总耗时。

## 示例压测命令

```bash
# 高负载：100 用户、5 分钟、高频弹幕送礼
go run main.go \
  -host=localhost:8080 \
  -room_id=1001 \
  -user_count=100 \
  -duration_seconds=300 \
  -chat_interval_ms=200 \
  -gift_interval_ms=800
```

## 验证方法

### 在线人数

打开 `http://localhost:8080/api/room/state?room_id=1001` 或在页面上观察 `online_count`。bot 启动后在线人数应上升到设定值，bot 结束后回落到 0。

### 排行榜

打开 `http://localhost:8080/api/room/rank?room_id=1001` 或在页面右侧 RankBoard 观察。bot 送礼后排行榜会动态变化，积分累加。

### 限流效果

设置 `-chat_interval_ms=100`（每个用户每 100ms 发一条弹幕），每秒约 10 条，超过后端 5 条/秒限制。

观察 `http://localhost:8080/api/room/state?room_id=1001` 中 `limited_count` 持续增长。

前端 SystemLog 区会显示限流通知。

### 消息统计

`/api/room/state` 返回的 `chat_count` 和 `gift_count` 记录累计数据，压测后观察数值应远大于 0。

### 广播延迟

后端日志中每次广播（target_clients > 0）会打印结构化日志：

```
INFO broadcast finished room_id=1001 type=chat target_clients=50 dropped_clients=0 latency_us=530
INFO broadcast finished room_id=1001 type=online target_clients=20 dropped_clients=0 latency_us=120
```

| 字段 | 说明 |
|------|------|
| `type` | 消息类型：chat / gift / rank / online |
| `target_clients` | 本次广播目标连接数 |
| `dropped_clients` | 因 send buffer 满而丢弃的数量 |
| `latency_us` | 广播耗时（微秒） |

## 压测结果记录

1. 执行压测后，打开 `docs/benchmark-result.md` 填写实际数据
2. 把 bot 输出的统计结果填入 Result Summary 表格
3. 调用 `curl localhost:8080/api/room/state?room_id=1001` 记录 chat_count / gift_count / limited_count
4. 调用 `curl localhost:8080/api/room/rank?room_id=1001` 记录排行榜 TOP10
5. 从后端日志摘取 `broadcast finished` 行，记录 target_clients / dropped_clients / latency_us
