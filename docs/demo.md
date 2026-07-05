# 功能演示验证指南

## 1. 初始化 MySQL

```bash
mysql -uroot -p < backend/sql/001_init_records.sql
```

## 2. 启动服务

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

## 3. 完整用户演示流程

### 3.1 打开房间大厅

浏览器打开 `http://localhost:3000`，自动跳转到 `/rooms` 房间大厅页面。

默认房间 `1001` 已在后端启动时自动创建，显示在列表中。

### 3.2 创建房间

点击「创建房间」按钮，输入房间标题（如 "今晚一起聊天"），点击创建。创建成功后自动跳转到该房间。

也可以通过 curl 创建：

```bash
curl -X POST "http://localhost:8080/api/rooms" \
  -H "Content-Type: application/json" \
  -d '{"title":"测试直播间","owner_name":"demo"}'
```

### 3.3 进入房间

在大厅点击房间卡片的「进入」按钮，或通过 URL 直接进入：
`http://localhost:3000/room/1001`

进入后自动连接 WebSocket，可以看到直播画面、弹幕区、礼物面板、排行榜等。

### 3.4 发送弹幕

在底部输入框输入弹幕内容，按回车或点击发送。弹幕会出现在左侧弹幕区。

### 3.5 送礼物

点击底部礼物按钮（小心心 / 火箭）。赠送后排行榜更新，系统日志显示送礼信息。

### 3.6 查看排行榜

右侧排行榜面板实时显示当前房间 TOP10 用户积分。

### 3.7 开两个窗口观察在线人数

再打开一个浏览器窗口（或隐身窗口）进入同一个房间，观察右侧 RoomStats 面板的在线人数变为 2。

### 3.8 启动 bot 压测

```bash
cd backend/bot
go run main.go -room_id=1001 -user_count=20 -duration_seconds=30
```

### 3.9 回到大厅观察统计

访问 `http://localhost:3000/rooms`，观察房间卡片的 `chat_count`、`gift_count`、`online_count` 数据。

## 4. 验证弹幕持久化

Bot 运行结束后：

```bash
curl "http://localhost:8080/api/room/chats?room_id=1001&limit=20"
```

## 5. 验证礼物持久化

```bash
curl "http://localhost:8080/api/room/gifts?room_id=1001&limit=20"
```

## 6. 验证房间 API

```bash
# 创建房间
curl -X POST "http://localhost:8080/api/rooms" \
  -H "Content-Type: application/json" \
  -d '{"title":"弹幕测试","owner_name":"test"}'

# 查看房间列表
curl "http://localhost:8080/api/rooms?limit=10"

# 查看单个房间
curl "http://localhost:8080/api/rooms/1001"

# 关闭房间
curl -X POST "http://localhost:8080/api/rooms/1001/close"

# 关闭后尝试进入（应返回 room is closed）
```

## 7. 验证全链路

```bash
# 运行 bot 压测
cd backend/bot
go run main.go -user_count=50 -duration_seconds=60

# 查询房间状态
curl "http://localhost:8080/api/room/state?room_id=1001"

# 查询弹幕流水
curl "http://localhost:8080/api/room/chats?room_id=1001&limit=100"

# 查询礼物流水
curl "http://localhost:8080/api/room/gifts?room_id=1001&limit=100"

# 查看房间详情（含 counts）
curl "http://localhost:8080/api/rooms/1001"
```
