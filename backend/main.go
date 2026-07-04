package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"liveroom-battle/common"
	"liveroom-battle/config"
	"liveroom-battle/controller"
	"liveroom-battle/dao"
	"liveroom-battle/hub"
	"liveroom-battle/model"
	"liveroom-battle/router"
	"liveroom-battle/service"
)

func main() {
	common.InitLogger()

	if err := config.Load("config/config.toml"); err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}

	rdb := common.InitRedis(config.Cfg.Redis)
	redisDao := dao.NewRedisDao(rdb)

	mysqlDB := common.InitMySQL(config.Cfg.MySQL)
	if err := mysqlDB.AutoMigrate(
		&model.Room{},
		&model.ChatMessage{},
		&model.GiftRecord{},
	); err != nil {
		slog.Error("auto migrate failed", "err", err)
		os.Exit(1)
	}

	roomDao := dao.NewRoomDao(mysqlDB)
	chatMessageDao := dao.NewChatMessageDao(mysqlDB)
	giftRecordDao := dao.NewGiftRecordDao(mysqlDB)

	roomHub := hub.NewRoomHub()

	roomSvc := service.NewRoomService(redisDao, roomHub, roomDao)
	rateLimitSvc := service.NewRateLimitService(redisDao, 5)
	rankSvc := service.NewRankService(redisDao)
	chatSvc := service.NewChatService(rateLimitSvc, roomHub, chatMessageDao)
	giftSvc := service.NewGiftService(redisDao, roomHub, giftRecordDao)

	dispatcher := service.NewMessageDispatcher()
	dispatcher.Register("chat", func(ctx context.Context, client *model.Client, msg *model.Message) {
		chatSvc.HandleChat(ctx, client, msg)
	})
	dispatcher.Register("gift", func(ctx context.Context, client *model.Client, msg *model.Message) {
		giftSvc.HandleGift(ctx, client, msg)
	})

	wsCtrl := controller.NewWSController(dispatcher, roomSvc)
	roomCtrl := controller.NewRoomController(roomSvc, rankSvc, chatMessageDao, giftRecordDao)

	r := router.Setup(wsCtrl, roomCtrl)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		addr := fmt.Sprintf(":%d", config.Cfg.Server.Port)
		slog.Info(fmt.Sprintf("server starting on %s", addr))
		if err := r.Run(addr); err != nil {
			slog.Error("server run failed", "err", err)
			os.Exit(1)
		}
	}()

	<-quit
	slog.Info("shutting down server...")
}
