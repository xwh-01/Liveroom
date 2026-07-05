# 功能演示验证指南

## 1. 启动 Redis

```bash
redis-server
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

浏览器打开 `http://localhost:3000`，自动跳转到 `/rooms` 房间大厅页面。

可以看到 3 个预置直播间卡片，分别展示：
- 房间标题、主播名
- 在线人数、弹幕数、礼物数
- 直播状态

### 4.2 进入直播间

点击任意房间的「进入直播间」按钮，跳转到 `/room/:roomId`。

进入后：
- 自动连接 WebSocket
- 顶部显示房间标题和主播名
- 左侧显示直播画面和弹幕区
- 右侧显示在线人数、排行榜、系统日志

### 4.3 发送弹幕

在底部输入框输入弹幕内容，按回车或点击发送。弹幕出现在左侧弹幕区。

### 4.4 送礼物

点击底部礼物按钮（小心心 / 火箭）。排行榜实时更新，系统日志显示送礼信息。

### 4.5 查看排行榜

右侧 RankBoard 展示当前房间 TOP10 用户积分。

### 4.6 开两个窗口观察在线人数

再打开一个浏览器窗口进入同一房间，观察：
- 右侧 RoomStats 在线人数变为 2
- 回到大厅看到 online_count 更新

### 4.7 启动 bot 压测

```bash
cd backend/bot
go run main.go -room_id=1001 -user_count=20 -duration_seconds=30
```

### 4.8 回到大厅观察统计

访问 `http://localhost:3000/rooms`，观察房间卡片的 chat_count、gift_count、online_count 数据变化。

## 5. API 验证

```bash
# 查看房间列表
curl "http://localhost:8080/api/rooms?limit=10"

# 查看单个房间
curl "http://localhost:8080/api/rooms/1001"

# 查看房间状态
curl "http://localhost:8080/api/room/state?room_id=1001"

# 查看排行榜
curl "http://localhost:8080/api/room/rank?room_id=1001"

# 关闭房间（开发环境）
curl -X POST "http://localhost:8080/api/admin/rooms/1001/close"

# 关闭后尝试 WebSocket 连接（应返回 room is closed）
```
