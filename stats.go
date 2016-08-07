package lnsq

import (
	"sync"
	"sync/atomic"
)

type CountStats struct {
	sync.RWMutex
	Count    int64
	SubStats map[string]*CountStats
}

func NewStats() *CountStats {
	return &CountStats{
		SubStats: make(map[string]*CountStats),
	}
}

func (s *CountStats) IncrCount() {
	atomic.AddInt64(&s.Count, 1)
}

func (s *CountStats) NewSubStats(name string) *CountStats {
	s.Lock()
	defer s.Unlock()

	if name == "" {
		name = randStatID()
	}
	countStats := NewStats()
	s.SubStats[name] = countStats
	return countStats
}
