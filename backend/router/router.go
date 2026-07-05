package router

import (
	"liveroom-battle/controller"
	"liveroom-battle/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(wsCtrl *controller.WSController, roomCtrl *controller.RoomController) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Cors())

	r.GET("/ws", wsCtrl.HandleWS)

	api := r.Group("/api")
	{
		api.GET("/room/state", roomCtrl.GetRoomState)
		api.GET("/room/rank", roomCtrl.GetRank)
		api.GET("/room/chats", roomCtrl.ListRecentChats)
		api.GET("/room/gifts", roomCtrl.ListRecentGifts)
		api.GET("/room/persist/state", roomCtrl.GetPersistState)
		api.GET("/room/pk/state", roomCtrl.GetPKState)
		api.GET("/room/pk/rank", roomCtrl.GetPKRank)
		api.GET("/rooms", roomCtrl.ListRooms)
		api.GET("/rooms/:room_id", roomCtrl.GetRoom)
	}

	admin := r.Group("/api/admin")
	{
		admin.POST("/rooms/:room_id/close", roomCtrl.CloseRoom)
		admin.POST("/room/pk/start", roomCtrl.StartPK)
		admin.POST("/room/pk/end", roomCtrl.EndPK)
	}

	return r
}
