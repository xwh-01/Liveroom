package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"liveroom-battle/dao"
	"liveroom-battle/hub"
	"liveroom-battle/model"
)

type PKService struct {
	redisDao *dao.RedisDao
	hub      *hub.RoomHub
}

func NewPKService(redisDao *dao.RedisDao, hub *hub.RoomHub) *PKService {
	return &PKService{redisDao: redisDao, hub: hub}
}

func (s *PKService) StartPK(ctx context.Context, roomID string, durationSeconds int) (*model.PKState, error) {
	existing, _ := s.redisDao.GetPKState(ctx, roomID)
	if existing != nil && existing.Status == "running" {
		return nil, errors.New("pk already running")
	}

	if durationSeconds <= 0 {
		durationSeconds = 300
	}

	now := time.Now()
	state := model.PKState{
		RoomID:    roomID,
		Status:    "running",
		RedScore:  0,
		BlueScore: 0,
		RedUsers:  0,
		BlueUsers: 0,
		Winner:    "",
		StartTime: strconv.FormatInt(now.Unix(), 10),
		EndTime:   strconv.FormatInt(now.Add(time.Duration(durationSeconds)*time.Second).Unix(), 10),
	}
	if err := s.redisDao.SetPKState(ctx, state); err != nil {
		return nil, err
	}
	s.broadcastPKState(ctx, roomID)
	return &state, nil
}

func (s *PKService) EndPK(ctx context.Context, roomID string) (*model.PKState, error) {
	state, err := s.redisDao.GetPKState(ctx, roomID)
	if err != nil || state == nil {
		return nil, errors.New("pk not found")
	}
	if state.Status != "running" {
		return nil, errors.New("pk not running")
	}

	redScore, blueScore, _ := s.redisDao.GetPKScore(ctx, roomID)
	state.Status = "ended"
	state.RedScore = redScore
	state.BlueScore = blueScore
	if redScore > blueScore {
		state.Winner = "red"
	} else if blueScore > redScore {
		state.Winner = "blue"
	} else {
		state.Winner = "draw"
	}

	redUsers, _ := s.redisDao.GetTeamUserCount(ctx, roomID, "red")
	blueUsers, _ := s.redisDao.GetTeamUserCount(ctx, roomID, "blue")
	state.RedUsers = redUsers
	state.BlueUsers = blueUsers

	if err := s.redisDao.SetPKState(ctx, *state); err != nil {
		return nil, err
	}
	s.broadcastPKState(ctx, roomID)
	return state, nil
}

func (s *PKService) JoinTeam(ctx context.Context, roomID, userID, team string) (*model.PKState, error) {
	if team != "red" && team != "blue" {
		return nil, errors.New("team must be red or blue")
	}

	existing, err := s.redisDao.GetUserTeam(ctx, roomID, userID)
	if err != nil {
		return nil, err
	}
	if existing != "" {
		return nil, fmt.Errorf("already in %s team", existing)
	}

	state, err := s.redisDao.GetPKState(ctx, roomID)
	if err != nil || state == nil || state.Status != "running" {
		return nil, errors.New("pk not running")
	}

	if err := s.redisDao.SetUserTeam(ctx, roomID, userID, team); err != nil {
		return nil, err
	}
	if err := s.redisDao.AddTeamUser(ctx, roomID, team, userID); err != nil {
		return nil, err
	}

	return s.GetPKState(ctx, roomID)
}

func (s *PKService) AddScore(ctx context.Context, roomID, userID string, score int) (*model.PKState, error) {
	team, err := s.redisDao.GetUserTeam(ctx, roomID, userID)
	if err != nil || team == "" {
		return nil, errors.New("user not in any team")
	}

	state, err := s.redisDao.GetPKState(ctx, roomID)
	if err != nil || state == nil || state.Status != "running" {
		return nil, errors.New("pk not running")
	}

	if err := s.redisDao.AddPKScore(ctx, roomID, userID, team, score); err != nil {
		return nil, err
	}

	return s.GetPKState(ctx, roomID)
}

func (s *PKService) GetPKState(ctx context.Context, roomID string) (*model.PKState, error) {
	state, err := s.redisDao.GetPKState(ctx, roomID)
	if err != nil || state == nil {
		return state, err
	}

	redScore, blueScore, _ := s.redisDao.GetPKScore(ctx, roomID)
	state.RedScore = redScore
	state.BlueScore = blueScore

	redUsers, _ := s.redisDao.GetTeamUserCount(ctx, roomID, "red")
	blueUsers, _ := s.redisDao.GetTeamUserCount(ctx, roomID, "blue")
	state.RedUsers = redUsers
	state.BlueUsers = blueUsers

	if state.Status == "running" && state.EndTime != "" {
		endTS, _ := strconv.ParseInt(state.EndTime, 10, 64)
		remaining := int(endTS - time.Now().Unix())
		if remaining < 0 {
			remaining = 0
		}
		state.RemainingSeconds = remaining
	}

	return state, nil
}

func (s *PKService) broadcastPKState(ctx context.Context, roomID string) {
	state, err := s.GetPKState(ctx, roomID)
	if err != nil || state == nil {
		return
	}
	payload, _ := model.NewResponse("pk_state", roomID, "", state)
	s.hub.Broadcast(roomID, "pk_state", payload)
}

func (s *PKService) GetTeamRank(ctx context.Context, roomID, team string, topN int) ([]model.RankItem, error) {
	return s.redisDao.GetTeamRank(ctx, roomID, team, topN)
}
