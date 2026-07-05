package controller

import (
	"net/http"
	"strconv"

	"liveroom-battle/model"
	"liveroom-battle/service"
	"liveroom-battle/utils"

	"github.com/gin-gonic/gin"
)

type PKController struct {
	pkSvc *service.PKService
}

func NewPKController(pkSvc *service.PKService) *PKController {
	return &PKController{pkSvc: pkSvc}
}

func (c *PKController) GetPKState(ctx *gin.Context) {
	roomID := ctx.Query("room_id")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	state, err := c.pkSvc.GetPKState(ctx.Request.Context(), roomID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrInternal)
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(state))
}

func (c *PKController) StartPK(ctx *gin.Context) {
	var req struct {
		RoomID          string `json:"room_id"`
		DurationSeconds int    `json:"duration_seconds"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	state, err := c.pkSvc.StartPK(ctx.Request.Context(), req.RoomID, req.DurationSeconds)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{Code: 400, Msg: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(state))
}

func (c *PKController) EndPK(ctx *gin.Context) {
	var req struct {
		RoomID string `json:"room_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	state, err := c.pkSvc.EndPK(ctx.Request.Context(), req.RoomID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.Response{Code: 400, Msg: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, utils.Success(state))
}

func (c *PKController) GetPKRank(ctx *gin.Context) {
	roomID := ctx.Query("room_id")
	team := ctx.Query("team")
	if roomID == "" || team == "" {
		ctx.JSON(http.StatusBadRequest, utils.ErrBadRequest)
		return
	}
	topN := 10
	if n := ctx.Query("top_n"); n != "" {
		if parsed, err := strconv.Atoi(n); err == nil {
			topN = parsed
		}
	}
	if topN > 100 {
		topN = 100
	}
	items, err := c.pkSvc.GetTeamRank(ctx.Request.Context(), roomID, team, topN)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrInternal)
		return
	}
	if items == nil {
		items = []model.RankItem{}
	}
	ctx.JSON(http.StatusOK, utils.Success(items))
}
