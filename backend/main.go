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

	db, err := common.InitMySQL(config.Cfg.MySQL)
	if err != nil {
		slog.Error("failed to init mysql", "err", err)
		os.Exit(1)
	}
	recordDao := dao.NewRecordDao(db)

	roomHub := hub.NewRoomHub()

	persistSvc := service.NewPersistService(recordDao, 10000)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	persistSvc.Start(ctx, 2)

	roomSvc := service.NewRoomService(redisDao, roomHub, persistSvc)
	roomManageSvc := service.NewRoomManageService(redisDao, roomHub)
	rateLimitSvc := service.NewRateLimitService(redisDao, 5)
	rankSvc := service.NewRankService(redisDao)
	chatSvc := service.NewChatService(rateLimitSvc, roomHub, redisDao, persistSvc)
	giftSvc := service.NewGiftService(redisDao, roomHub, persistSvc)

	dispatcher := service.NewMessageDispatcher()
	dispatcher.Register("chat", func(ctx context.Context, client *model.Client, msg *model.Message) {
		chatSvc.HandleChat(ctx, client, msg)
	})
	dispatcher.Register("gift", func(ctx context.Context, client *model.Client, msg *model.Message) {
		giftSvc.HandleGift(ctx, client, msg)
	})

	if err := roomManageSvc.EnsureDefaultRooms(context.Background()); err != nil {
		slog.Error("failed to ensure default rooms", "err", err)
		os.Exit(1)
	}

	wsCtrl := controller.NewWSController(dispatcher, roomSvc, roomManageSvc)
	roomCtrl := controller.NewRoomController(roomSvc, roomManageSvc, rankSvc, recordDao)

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
