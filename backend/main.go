package main

import (
	"context"
	"encoding/json"
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

func handleJoinTeam(ctx context.Context, client *model.Client, msg *model.Message, pkSvc *service.PKService) {
	var teamData model.TeamData
	if err := json.Unmarshal(msg.Data, &teamData); err != nil {
		slog.Error("unmarshal join_team failed", "err", err)
		return
	}
	state, err := pkSvc.JoinTeam(ctx, msg.RoomID, msg.UserID, teamData.Team)
	if err != nil {
		reply, _ := model.NewResponse("system", msg.RoomID, "", model.SystemData{
			Content: err.Error(),
		})
		client.Send <- reply
		return
	}
	systemMsg := ""
	if teamData.Team == "red" {
		systemMsg = "你已加入红队"
	} else {
		systemMsg = "你已加入蓝队"
	}
	reply, _ := model.NewResponse("system", msg.RoomID, "", model.SystemData{Content: systemMsg})
	client.Send <- reply

	if state != nil {
		payload, _ := model.NewResponse("pk_state", msg.RoomID, "", state)
		pkSvc.GetPKState(ctx, msg.RoomID)
		client.Send <- payload
	}
}

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

	roomSvc := service.NewRoomService(redisDao, roomHub)
	roomManageSvc := service.NewRoomManageService(redisDao, roomHub)
	pkSvc := service.NewPKService(redisDao, roomHub)
	rateLimitSvc := service.NewRateLimitService(redisDao, 5)
	rankSvc := service.NewRankService(redisDao)
	chatSvc := service.NewChatService(rateLimitSvc, roomHub, redisDao, persistSvc)
	giftSvc := service.NewGiftService(redisDao, roomHub, persistSvc, pkSvc)

	dispatcher := service.NewMessageDispatcher()
	dispatcher.Register("chat", func(ctx context.Context, client *model.Client, msg *model.Message) {
		chatSvc.HandleChat(ctx, client, msg)
	})
	dispatcher.Register("gift", func(ctx context.Context, client *model.Client, msg *model.Message) {
		giftSvc.HandleGift(ctx, client, msg)
	})
	dispatcher.Register("join_team", func(ctx context.Context, client *model.Client, msg *model.Message) {
		handleJoinTeam(ctx, client, msg, pkSvc)
	})

	if err := roomManageSvc.EnsureDefaultRooms(context.Background()); err != nil {
		slog.Error("failed to ensure default rooms", "err", err)
		os.Exit(1)
	}

	wsCtrl := controller.NewWSController(dispatcher, roomSvc, roomManageSvc)
	roomCtrl := controller.NewRoomController(roomSvc, roomManageSvc, rankSvc, recordDao, persistSvc, pkSvc)

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
