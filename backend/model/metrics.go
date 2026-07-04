package model

import (
	"sync"
	"sync/atomic"
)

type RoomMetrics struct {
	ChatCount atomic.Int64
	GiftCount atomic.Int64
}

type MetricsStore struct {
	mu      sync.RWMutex
	metrics map[string]*RoomMetrics
}

func NewMetricsStore() *MetricsStore {
	return &MetricsStore{metrics: make(map[string]*RoomMetrics)}
}

func (s *MetricsStore) GetOrCreate(roomID string) *RoomMetrics {
	s.mu.RLock()
	m, ok := s.metrics[roomID]
	s.mu.RUnlock()
	if ok {
		return m
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok = s.metrics[roomID]; ok {
		return m
	}
	m = &RoomMetrics{}
	s.metrics[roomID] = m
	return m
}

func (s *MetricsStore) GetChatCount(roomID string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if m, ok := s.metrics[roomID]; ok {
		return m.ChatCount.Load()
	}
	return 0
}

func (s *MetricsStore) GetGiftCount(roomID string) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if m, ok := s.metrics[roomID]; ok {
		return m.GiftCount.Load()
	}
	return 0
}
