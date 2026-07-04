package service

import (
	"context"

	"liveroom-battle/dao"
	"liveroom-battle/model"
)

type RankService struct {
	dao *dao.RedisDao
}

func NewRankService(dao *dao.RedisDao) *RankService {
	return &RankService{dao: dao}
}

func (s *RankService) GetTop10(ctx context.Context, roomID string) ([]model.RankItem, error) {
	return s.dao.GetTopRank(ctx, roomID, 10)
}
