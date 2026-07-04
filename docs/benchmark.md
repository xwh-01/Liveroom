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

#### 自定义压测参数

```bash
go run main.go \
  -room_id=1001 \
  -user_count=50 \
  -duration_seconds=120 \
  -chat_interval_ms=300 \
  -gift_interval_ms=1000
```

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-room_id` | `1001` | 目标房间 ID |
| `-user_count` | `20` | 模拟用户数 |
| `-duration_seconds` | `60` | 压测持续时间（秒） |
| `-chat_interval_ms` | `500` | 每个用户发弹幕间隔（毫秒） |
| `-gift_interval_ms` | `1500` | 每个用户送礼间隔（毫秒） |

## 示例压测命令

```bash
# 高负载：100 用户、5 分钟、高频弹幕送礼
go run main.go \
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

前端 `<--SystemLog-->` 区会显示限流通知。

### 消息统计

`/api/room/state` 返回的 `chat_count` 和 `gift_count` 记录累计数据，压测后观察数值应远大于 0。

### 广播延迟

后端日志中会打印每次广播的延迟（超过 10ms 时）：

```
INFO broadcast done room_id=1001 clients=50 dropped=0 latency_us=1234
```
