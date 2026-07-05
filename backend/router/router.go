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
	}

	return r
}
