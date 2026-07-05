# 功能演示验证指南

## 1. 初始化 MySQL

```bash
mysql -uroot -p < backend/sql/001_init_records.sql
```

这会创建 `liveroom` 数据库以及 `chat_records` 和 `gift_records` 两张表。

## 2. 修改 MySQL 连接配置

编辑 `backend/config/config.toml`，将 `[mysql]` 中的 `dsn` 修改为你的 MySQL 账号密码：

```toml
[mysql]
dsn = "root:your_password@tcp(127.0.0.1:3306)/liveroom?charset=utf8mb4&parseTime=True&loc=Local"
max_open_conns = 20
max_idle_conns = 10
conn_max_lifetime_seconds = 300
```

## 3. 启动服务

按顺序启动：

```bash
# 终端1: 启动 Redis
redis-server

# 终端2: 启动后端
cd backend
go run main.go

# 终端3: 启动前端
cd vue-frontend
npm install
npm run dev
```

## 4. 验证弹幕持久化

### 方式一：通过前端页面

1. 浏览器打开 `http://localhost:3000`
2. 输入房间号（如 `1001`）和用户名，点击进入
3. 在弹幕输入框发送一条弹幕
4. 用 curl 查询弹幕记录：

```bash
curl "http://localhost:8080/api/room/chats?room_id=1001&limit=20"
```

### 方式二：直接查询数据库

```sql
SELECT * FROM chat_records WHERE room_id = '1001' ORDER BY created_at DESC LIMIT 20;
```

### 方式三：通过 Bot 批量验证

```bash
cd backend/bot
go run main.go -duration_seconds=10
```

Bot 运行结束后，查询弹幕记录：

```bash
curl "http://localhost:8080/api/room/chats?room_id=1001&limit=100"
```

## 5. 验证礼物持久化

```bash
curl "http://localhost:8080/api/room/gifts?room_id=1001&limit=20"
```

直接查询数据库：

```sql
SELECT * FROM gift_records WHERE room_id = '1001' ORDER BY created_at DESC LIMIT 20;
```

## 6. 验证落库失败不影响广播

即使 MySQL 不可用（关闭 MySQL 或连接断开），后端仍应正常运行：
- 弹幕和礼物广播正常
- Worker 打印 error 日志但不 panic
- `/api/room/state` 返回正常数据

## 7. 验证队列满丢弃

可以通过调小 queueSize 来模拟（临时修改 main.go 中 `NewPersistService(recordDao, 10)`），然后用 bot 高压写入。

观察日志中是否出现：

```
WARN persist queue full, event dropped type=chat room_id=1001 user_id=...
```

同时 `/api/room/state` 中的 `persist_dropped_count` 会增长：

```bash
curl "http://localhost:8080/api/room/state?room_id=1001"
```

返回结果中包含：

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "room_id": "1001",
    "online_count": 5,
    "limited_count": 0,
    "chat_count": 1234,
    "gift_count": 567,
    "persist_dropped_count": 10
  }
}
```

## 8. 验证 API 限制

limit 参数验证：
- `limit` 不传默认 20
- `limit` 传超过 100 会被截断为 100
- `room_id` 不传返回 `{"code":400,"msg":"bad request"}`

```bash
# 默认 20 条
curl "http://localhost:8080/api/room/chats?room_id=1001"

# 指定 50 条
curl "http://localhost:8080/api/room/chats?room_id=1001&limit=50"

# 超过 100 截断为 100
curl "http://localhost:8080/api/room/chats?room_id=1001&limit=999"

# 缺少 room_id
curl "http://localhost:8080/api/room/chats"
```

## 9. 压测后验证全链路

```bash
# 运行 bot 压测 60 秒
cd backend/bot
go run main.go -user_count=50 -duration_seconds=60

# 查询房间状态
curl "http://localhost:8080/api/room/state?room_id=1001"

# 查询弹幕流水
curl "http://localhost:8080/api/room/chats?room_id=1001&limit=100"

# 查询礼物流水
curl "http://localhost:8080/api/room/gifts?room_id=1001&limit=100"

# 检查 persist_dropped_count
curl -s "http://localhost:8080/api/room/state?room_id=1001" | grep persist_dropped_count
```
