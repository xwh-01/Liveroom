package controller

import (
	"net/http"
	"strconv"

	"liveroom-battle/dao"
	"liveroom-battle/model"
	"liveroom-battle/service"
	"liveroom-battle/utils"

	"github.com/gin-gonic/gin"
)

type RoomController struct {
	roomSvc       *service.RoomService
	roomManageSvc *service.RoomManageService
	rankSvc       *service.RankService
	recordDao     *dao.RecordDao
}

func NewRoomController(roomSvc *service.RoomService, roomManageSvc *service.RoomManageService, rankSvc *service.RankService, recordDao *dao.RecordDao) *RoomController {
	return &RoomController{
		roomSvc:       roomSvc,
		roomManageSvc: roomManageSvc,
		rankSvc:       rankSvc,
		recordDao:     recordDao,
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

func (c *RoomController) ListRooms(ctx *gin.Context) {
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
	rooms, err := c.roomManageSvc.ListLiveRooms(ctx.Request.Context(), limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrInternal)
		return
	}
	if rooms == nil {
		rooms = []model.RoomMeta{}
	}
	ctx.JSON(http.StatusOK, utils.Success(rooms))
}

func (c *RoomController) GetRoom(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	meta, err := c.roomManageSvc.GetRoom(ctx.Request.Context(), roomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrInternal)
		return
	}
	if meta == nil {
		ctx.JSON(http.StatusNotFound, utils.Response{Code: 404, Msg: "room not found"})
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(meta))
}

func (c *RoomController) CloseRoom(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	if err := c.roomManageSvc.CloseRoom(ctx.Request.Context(), roomID); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{Code: 400, Msg: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(nil))
}
