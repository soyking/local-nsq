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
}

func NewLocalNSQ() *LocalNSQ {
	return &LocalNSQ{
		channels: make(map[string]*Channel),
	}
}

func (l *LocalNSQ) Subscribe(channel, topic string, callback Callback, maxInFlight, concurrency int) {
	l.channelsLock.Lock()
	defer l.channelsLock.Unlock()

	c, ok := l.channels[channel]
	if !ok {
		c = NewChannel(channel)
		l.channels[channel] = c
	}

	c.Subscribe(topic, callback, maxInFlight, concurrency)
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
	l.channelsLock.RLock()
	defer l.channelsLock.RUnlock()

	channelsStats := make(map[string]int64)
	for name, channel := range l.channels {
		channelsStats[name] = channel.Stat.Count()
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
