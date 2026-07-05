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
	roomSvc   *service.RoomService
	rankSvc   *service.RankService
	recordDao *dao.RecordDao
}

func NewRoomController(roomSvc *service.RoomService, rankSvc *service.RankService, recordDao *dao.RecordDao) *RoomController {
	return &RoomController{
		roomSvc:   roomSvc,
		rankSvc:   rankSvc,
		recordDao: recordDao,
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

func (c *RoomController) ListRecentChats(ctx *gin.Context) {
	roomID := ctx.Query("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	limit := 20
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	records, err := c.recordDao.ListRecentChatRecords(ctx.Request.Context(), roomID, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrInternal)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(records))
}

func (c *RoomController) ListRecentGifts(ctx *gin.Context) {
	roomID := ctx.Query("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	limit := 20
	if l := ctx.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	records, err := c.recordDao.ListRecentGiftRecords(ctx.Request.Context(), roomID, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrInternal)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(records))
}
