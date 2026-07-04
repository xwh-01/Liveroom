package controller

import (
	"net/http"
	"strconv"

	"liveroom-battle/dao"
	"liveroom-battle/service"
	"liveroom-battle/utils"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	roomSvc        *service.RoomService
	rankSvc        *service.RankService
	chatMessageDao *dao.ChatMessageDao
	giftRecordDao  *dao.GiftRecordDao
}

func NewRoomController(roomSvc *service.RoomService, rankSvc *service.RankService, chatMessageDao *dao.ChatMessageDao, giftRecordDao *dao.GiftRecordDao) *RoomController {
	return &RoomController{
		roomSvc:        roomSvc,
		rankSvc:        rankSvc,
		chatMessageDao: chatMessageDao,
		giftRecordDao:  giftRecordDao,
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

func (c *RoomController) GetChatHistory(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	limitStr := ctx.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	messages, err := c.chatMessageDao.ListRecent(ctx.Request.Context(), roomID, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrInternal)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(messages))
}

func (c *RoomController) GetGiftHistory(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	limitStr := ctx.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}
	records, err := c.giftRecordDao.ListRecent(ctx.Request.Context(), roomID, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrInternal)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(records))
}
