# Benchmark Result

## Environment

- OS: 待填写
- CPU: 待填写
- Memory: 待填写
- Go Version: 待填写
- Redis Version: 待填写
- Backend Mode: single instance
- Test Date: 待填写

## Test Command

```bash
cd backend/bot
go run . \
  -host=localhost:8080 \
  -room_id=1001 \
  -user_count=100 \
  -duration_seconds=120 \
  -chat_interval_ms=300 \
  -gift_interval_ms=1000
```

## Result Summary

| Metric | Value |
|--------|-------|
| Conn Success | 待填写 |
| Conn Errors | 待填写 |
| Chat Sent | 待填写 |
| Gift Sent | 待填写 |
| Send Errors | 待填写 |
| Disconnects | 待填写 |
| Online Count Peak | 待填写 |
| Chat Count (Redis) | 待填写 |
| Gift Count (Redis) | 待填写 |
| Limited Count (Redis) | 待填写 |

## Room State Sample

```bash
curl "http://localhost:8080/api/room/state?room_id=1001"
```

```json
{
  "code": 0,
  "msg": "ok",
  "data": {
    "room_id": "1001",
    "online_count": 待填写,
    "limited_count": 待填写,
    "chat_count": 待填写,
    "gift_count": 待填写
  }
}
```

## Rank Sample

```bash
curl "http://localhost:8080/api/room/rank?room_id=1001"
```

```json
{
  "code": 0,
  "msg": "ok",
  "data": [
    { "user_id": "bot_1", "score": 待填写 },
    { "user_id": "bot_2", "score": 待填写 }
  ]
}
```

## Broadcast Log Sample

从后端日志中摘取 `broadcast finished` 日志（target_clients > 0 时打印）：

```
INFO broadcast finished room_id=1001 type=chat target_clients=待填写 dropped_clients=待填写 latency_us=待填写
INFO broadcast finished room_id=1001 type=gift target_clients=待填写 dropped_clients=待填写 latency_us=待填写
INFO broadcast finished room_id=1001 type=rank target_clients=待填写 dropped_clients=待填写 latency_us=待填写
INFO broadcast finished room_id=1001 type=online target_clients=待填写 dropped_clients=待填写 latency_us=待填写
```

## Conclusion

本次测试模拟了 N 个用户同时进入直播间，在持续 T 秒内发送聊天和礼物消息。系统能够维持 WebSocket 连接、完成消息广播、更新 Redis 礼物排行榜，并通过限流计数和广播日志观察高频消息下的系统行为。

> 注意：上表中所有"待填写"字段请在实际压测后替换为真实数据。
