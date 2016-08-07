package lnsq

import (
	"errors"
	"sync"
)

var (
	ErrNoSuchChannel = errors.New("no such channel")
)

type LocalNSQ struct {
	channelsLock sync.RWMutex
	channels     map[string]*Channel
	Stats        *CountStats
}

func NewLocalNSQ() *LocalNSQ {
	return &LocalNSQ{
		channels: make(map[string]*Channel),
		Stats:    NewStats(),
	}
}

func (l *LocalNSQ) Subscribe(channel, topic string, callback Callback, maxInFlight, concurrency int) {
	c := l.getChannel(channel)
	c.Subscribe(topic, callback, maxInFlight, concurrency)
}

func (l *LocalNSQ) getChannel(channel string) *Channel {
	l.channelsLock.Lock()
	defer l.channelsLock.Unlock()

	c, ok := l.channels[channel]
	if !ok {
		c = NewChannel(channel, l.Stats.NewSubStats(channel))
		l.channels[channel] = c
	}

	return c
}

func (l *LocalNSQ) Dispatch(channel string, msg interface{}) error {
	l.channelsLock.RLock()
	defer l.channelsLock.RUnlock()

	if c, ok := l.channels[channel]; !ok {
		return ErrNoSuchChannel
	} else {
		c.Dispatch(msg)
		return nil
	}
}

func (l *LocalNSQ) ChannelsStats() map[string]int64 {
	l.Stats.RLock()
	defer l.Stats.RUnlock()

	channelsStats := make(map[string]int64)
	for channel, stats := range l.Stats.SubStats {
		channelsStats[channel] = stats.Count
	}
	return channelsStats
}

func (l *LocalNSQ) TopicsStats(channel string) map[string]int64 {
	l.channelsLock.RLock()
	defer l.channelsLock.RUnlock()

	if c, ok := l.channels[channel]; !ok {
		return map[string]int64{}
	} else {
		return c.TopicsStats()
	}
}

func (l *LocalNSQ) CallbacksStats(channel, topic string) map[string]int64 {
	l.channelsLock.RLock()
	defer l.channelsLock.RUnlock()

	if c, ok := l.channels[channel]; !ok {
		return map[string]int64{}
	} else {
		return c.CallbacksStats(topic)
	}
}
