# 功能演示验证指南

## 1. 启动基础设施

```bash
# 方式一：docker compose
docker compose up -d

# 方式二：手动启动
redis-server

# MySQL 初始化
mysql -uroot -p < backend/sql/001_init_records.sql
```

## 2. 启动后端

```bash
cd backend
go run main.go
```

后端启动时自动初始化 3 个预置直播间（1001 / 1002 / 1003）。

## 3. 启动前端

```bash
cd vue-frontend
npm install
npm run dev
```

## 4. 完整演示流程

### 4.1 打开房间大厅

浏览器打开 `http://localhost:3000`，自动跳转到 `/rooms`。

### 4.2 进入直播间

点击卡片进入房间，自动连接 WebSocket。

### 4.3 发送弹幕

输入弹幕内容，按回车发送。弹幕出现在左侧弹幕区。

### 4.4 送礼物

点击礼物按钮（小心心 / 火箭），排行榜实时更新。

### 4.5 查询弹幕流水

```bash
curl "http://localhost:8080/api/room/chats?room_id=1001&limit=10"
```

### 4.6 查询礼物流水

```bash
curl "http://localhost:8080/api/room/gifts?room_id=1001&limit=10"
```

### 4.7 查看持久化队列状态

```bash
curl "http://localhost:8080/api/room/persist/state"
```

### 4.8 bot 压测

```bash
cd backend/bot
go run main.go -room_id=1001 -user_count=20 -duration_seconds=30
```

## 5. 全链路验证

```bash
# 查看房间列表
curl "http://localhost:8080/api/rooms"

# 查看房间状态
curl "http://localhost:8080/api/room/state?room_id=1001"

# 查看排行榜
curl "http://localhost:8080/api/room/rank?room_id=1001"

# 发送弹幕后查询
curl "http://localhost:8080/api/room/chats?room_id=1001&limit=10"

# 送礼物后查询
curl "http://localhost:8080/api/room/gifts?room_id=1001&limit=10"

# 队列状态
curl "http://localhost:8080/api/room/persist/state"
```
