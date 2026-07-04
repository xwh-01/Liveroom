package controller

import (
	"net/http"

	"liveroom-battle/service"
	"liveroom-battle/utils"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	roomSvc *service.RoomService
	rankSvc *service.RankService
}

func NewRoomController(roomSvc *service.RoomService, rankSvc *service.RankService) *RoomController {
	return &RoomController{
		roomSvc: roomSvc,
		rankSvc: rankSvc,
	}
}

func (c *RoomController) GetRoomState(ctx *gin.Context) {
	roomID := ctx.Query("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	state := c.roomSvc.GetRoomState(ctx.Request.Context(), roomID)
	ctx.JSON(http.StatusOK, utils.Success(state))
}

func (c *RoomController) GetRank(ctx *gin.Context) {
	roomID := ctx.Query("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	rankings, err := c.rankSvc.GetTop10(ctx.Request.Context(), roomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrInternal)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(rankings))
}
