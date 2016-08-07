package lnsq

import (
	"sync"
	"sync/atomic"
)

type CountStats struct {
	sync.RWMutex
	count int64
	stats map[string]*CountStats
}

func NewStats() *CountStats {
	return &CountStats{
		stats: make(map[string]*CountStats),
	}
}

func (s *CountStats) Count() {
	atomic.AddInt64(&s.count, 1)
}

func (s *CountStats) NewSubStats(name string) *CountStats {
	s.Lock()
	defer s.Unlock()

	if name == "" {
		name = randStatID()
	}
	countStats := NewStats()
	s.stats[name] = countStats
	return countStats
}
