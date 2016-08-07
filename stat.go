package lnsq

import (
	"sync/atomic"
)

type CountStat struct {
	count int64
}

func NewCountStat() *CountStat {
	return &CountStat{}
}

func (s *CountStat) IncrCount() {
	atomic.AddInt64(&s.count, 1)
}

func (s *CountStat) Count() int64 {
	return atomic.LoadInt64(&s.count)
}
