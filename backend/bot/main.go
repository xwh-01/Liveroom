package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

var (
	host        = flag.String("host", "localhost:8080", "server host:port")
	roomID      = flag.String("room_id", "1001", "room ID")
	userCount   = flag.Int("user_count", 20, "number of simulated users")
	durationSec = flag.Int("duration_seconds", 60, "test duration in seconds")
	chatMs      = flag.Int("chat_interval_ms", 500, "chat interval in ms")
	giftMs      = flag.Int("gift_interval_ms", 1500, "gift interval in ms")
	enablePK    = flag.Bool("enable_pk", true, "enable red-blue PK")
	redRatio    = flag.Int("red_ratio", 50, "percentage of red team users")
)

type Stats struct {
	Connected   int64
	ConnErrors  int64
	ChatSent    int64
	GiftSent    int64
	SendErrors  int64
	Disconnects int64
	RedUsers    int64
	BlueUsers   int64
	RedScore    int64
	BlueScore   int64
}

type Message struct {
	Type   string      `json:"type"`
	RoomID string      `json:"room_id"`
	UserID string      `json:"user_id"`
	Data   interface{} `json:"data"`
}

var stats Stats

func main() {
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	duration := time.Duration(*durationSec) * time.Second

	fmt.Println("=== LiveRoom Battle Bot ===")
	fmt.Printf("Target:      ws://%s/ws\n", *host)
	fmt.Printf("Room:        %s\n", *roomID)
	fmt.Printf("Users:       %d\n", *userCount)
	fmt.Printf("Duration:    %s\n", duration)
	fmt.Printf("Chat  every: %dms\n", *chatMs)
	fmt.Printf("Gift  every: %dms\n", *giftMs)
	fmt.Printf("PK Mode:     %v (red_ratio=%d%%)\n", *enablePK, *redRatio)
	fmt.Println("===========================")

	startTime := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < *userCount; i++ {
		wg.Add(1)
		userID := fmt.Sprintf("bot_%d", i+1)
		go simulateUser(userID, duration, &wg)
		time.Sleep(30 * time.Millisecond)
	}

	wg.Wait()
	elapsed := time.Since(startTime)

	fmt.Println()
	fmt.Println("=== Results ===")
	fmt.Printf("连接成功:      %d\n", stats.Connected)
	fmt.Printf("连接失败:      %d\n", stats.ConnErrors)
	fmt.Printf("Chat 发送:     %d\n", stats.ChatSent)
	fmt.Printf("Gift 发送:     %d\n", stats.GiftSent)
	fmt.Printf("写入错误:      %d\n", stats.SendErrors)
	fmt.Printf("断开连接:      %d\n", stats.Disconnects)
	fmt.Printf("红队用户:      %d\n", stats.RedUsers)
	fmt.Printf("蓝队用户:      %d\n", stats.BlueUsers)
	fmt.Printf("总耗时:        %s\n", elapsed.Round(time.Millisecond))
	fmt.Println("==============")
}

func simulateUser(userID string, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()

	u := url.URL{
		Scheme: "ws",
		Host:   *host,
		Path:   "/ws",
	}
	q := u.Query()
	q.Set("room_id", *roomID)
	q.Set("user_id", userID)
	u.RawQuery = q.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		atomic.AddInt64(&stats.ConnErrors, 1)
		fmt.Printf("[%s] connect failed: %v\n", userID, err)
		return
	}
	atomic.AddInt64(&stats.Connected, 1)
	defer func() {
		conn.Close()
		atomic.AddInt64(&stats.Disconnects, 1)
	}()

	if *enablePK {
		team := "blue"
		if rand.Intn(100) < *redRatio {
			team = "red"
		}
		msg := Message{
			Type:   "join_team",
			RoomID: *roomID,
			UserID: userID,
			Data:   map[string]string{"team": team},
		}
		raw, _ := json.Marshal(msg)
		conn.WriteMessage(websocket.TextMessage, raw)
		if team == "red" {
			atomic.AddInt64(&stats.RedUsers, 1)
		} else {
			atomic.AddInt64(&stats.BlueUsers, 1)
		}
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	chatTicker := time.NewTicker(time.Duration(*chatMs) * time.Millisecond)
	giftTicker := time.NewTicker(time.Duration(*giftMs) * time.Millisecond)
	defer chatTicker.Stop()
	defer giftTicker.Stop()

	timeout := time.After(duration)

	for {
		select {
		case <-timeout:
			return
		case <-chatTicker.C:
			msg := buildChat(userID)
			raw, _ := json.Marshal(msg)
			if err := conn.WriteMessage(websocket.TextMessage, raw); err != nil {
				atomic.AddInt64(&stats.SendErrors, 1)
				return
			}
			atomic.AddInt64(&stats.ChatSent, 1)
		case <-giftTicker.C:
			msg := buildGift(userID)
			raw, _ := json.Marshal(msg)
			if err := conn.WriteMessage(websocket.TextMessage, raw); err != nil {
				atomic.AddInt64(&stats.SendErrors, 1)
				return
			}
			atomic.AddInt64(&stats.GiftSent, 1)
		case <-done:
			return
		}
	}
}

func buildChat(userID string) Message {
	return Message{
		Type:   "chat",
		RoomID: *roomID,
		UserID: userID,
		Data:   map[string]string{"content": randomChat()},
	}
}

func buildGift(userID string) Message {
	giftType := "heart"
	if rand.Float64() < 0.3 {
		giftType = "rocket"
	}
	return Message{
		Type:   "gift",
		RoomID: *roomID,
		UserID: userID,
		Data:   map[string]string{"gift_type": giftType},
	}
}

var chatMessages = []string{
	"主播好棒！",
	"666666",
	"来了来了",
	"支持主播",
	"好厉害！",
	"加油加油",
	"前排围观",
	"这个绝了",
	"哈哈哈哈",
	"太强了吧",
	"给主播点赞",
	"第一次来，请多关照",
	"节目开始了",
	"能不能点歌",
	"主播今天状态不错",
	"弹幕测试",
	"路过一下",
	"支持一波",
	"关注了",
	"礼物走起",
}

func randomChat() string {
	return chatMessages[rand.Intn(len(chatMessages))]
}
