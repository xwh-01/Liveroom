package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type Stats struct {
	ChatSent  int64
	GiftSent  int64
	Errors    int64
	Connected int64
}

type Message struct {
	Type   string      `json:"type"`
	RoomID string      `json:"room_id"`
	UserID string      `json:"user_id"`
	Data   interface{} `json:"data"`
}

var stats Stats

func main() {
	rand.Seed(time.Now().UnixNano())

	host := "localhost:8080"
	roomID := "1001"
	userCount := 20
	highFreqCount := 4
	duration := 60 * time.Second

	fmt.Println("=== LiveRoom Battle Bot ===")
	fmt.Printf("Target: ws://%s/ws\n", host)
	fmt.Printf("Room: %s\n", roomID)
	fmt.Printf("Users: %d (high-freq: %d)\n", userCount, highFreqCount)
	fmt.Printf("Duration: %s\n", duration)
	fmt.Println("===========================")

	var wg sync.WaitGroup

	for i := 0; i < userCount; i++ {
		wg.Add(1)
		userID := fmt.Sprintf("bot_%d", i+1)
		isHighFreq := i < highFreqCount
		go simulateUser(userID, roomID, host, duration, isHighFreq, &wg)
		time.Sleep(50 * time.Millisecond)
	}

	wg.Wait()

	fmt.Println()
	fmt.Println("=== Results ===")
	fmt.Printf("Chat Sent:  %d\n", stats.ChatSent)
	fmt.Printf("Gift Sent:  %d\n", stats.GiftSent)
	fmt.Printf("Errors:     %d\n", stats.Errors)
	fmt.Println("==============")
}

func simulateUser(userID, roomID, host string, duration time.Duration, highFreq bool, wg *sync.WaitGroup) {
	defer wg.Done()

	u := url.URL{
		Scheme: "ws",
		Host:   host,
		Path:   "/ws",
	}
	q := u.Query()
	q.Set("room_id", roomID)
	q.Set("user_id", userID)
	u.RawQuery = q.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		atomic.AddInt64(&stats.Errors, 1)
		fmt.Printf("[%s] connect failed: %v\n", userID, err)
		return
	}
	defer conn.Close()

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

	ticker := time.NewTicker(getInterval(highFreq))
	defer ticker.Stop()

	timeout := time.After(duration)

	for {
		select {
		case <-timeout:
			return
		case <-ticker.C:
			msg := buildMessage(userID, roomID)
			raw, _ := json.Marshal(msg)
			if err := conn.WriteMessage(websocket.TextMessage, raw); err != nil {
				atomic.AddInt64(&stats.Errors, 1)
				return
			}
		case <-done:
			return
		}
	}
}

func buildMessage(userID, roomID string) Message {
	r := rand.Float64()
	if r < 0.65 {
		atomic.AddInt64(&stats.ChatSent, 1)
		content := randomChat()
		return Message{
			Type:   "chat",
			RoomID: roomID,
			UserID: userID,
			Data: map[string]string{
				"content": content,
			},
		}
	}
	atomic.AddInt64(&stats.GiftSent, 1)
	giftType := "heart"
	if rand.Float64() < 0.3 {
		giftType = "rocket"
	}
	return Message{
		Type:   "gift",
		RoomID: roomID,
		UserID: userID,
		Data: map[string]string{
			"gift_type": giftType,
		},
	}
}

func getInterval(highFreq bool) time.Duration {
	if highFreq {
		return time.Duration(50+rand.Intn(150)) * time.Millisecond
	}
	return time.Duration(200+rand.Intn(1800)) * time.Millisecond
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
