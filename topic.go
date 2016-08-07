package lnsq

import (
	"sync"
)

type Topic struct {
	Name string
	Stat *CountStat
	ch   chan interface{}

	callbacksLock sync.RWMutex
	callbacks     map[string]*CallbackWrapper
}

func NewTopic(name string, maxInFlight int) *Topic {
	return &Topic{
		Name:      name,
		Stat:      NewCountStat(),
		ch:        make(chan interface{}, maxInFlight),
		callbacks: make(map[string]*CallbackWrapper),
	}
}

func (t *Topic) Subscribe(c Callback, concurrency int) {
	t.callbacksLock.Lock()
	defer t.callbacksLock.Unlock()

	for i := 0; i < concurrency; i++ {
		cw := NewCallbackWrapper(c)
		t.callbacks[randStatID()] = cw
		go t.callback(cw)
	}
}

func (t *Topic) callback(c *CallbackWrapper) {
	for msg := range t.ch {
		c.Handler(msg)
	}
}

func (t *Topic) Dispatch(msg interface{}) {
	t.Stat.IncrCount()
	t.ch <- msg
}

func (t *Topic) CallbacksStats() map[string]int64 {
	t.callbacksLock.RLock()
	defer t.callbacksLock.RUnlock()

	callbacksStats := make(map[string]int64)
	for name, callback := range t.callbacks {
		callbacksStats[name] = callback.Stat.Count()
	}
	return callbacksStats
}
