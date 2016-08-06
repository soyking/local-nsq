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
	channels     map[string]*channel
}

func NewLocalNSQ() *LocalNSQ {
	return &LocalNSQ{channels: make(map[string]*channel)}
}

func (l *LocalNSQ) Subscribe(channel, topic string, callback Callback, maxInFlight, concurrency int) {
	c := l.getChannel(channel)
	c.subscribe(topic, callback, maxInFlight, concurrency)
}

func (l *LocalNSQ) getChannel(channel string) *channel {
	l.channelsLock.Lock()
	defer l.channelsLock.Unlock()

	c, ok := l.channels[channel]
	if !ok {
		c = newChannel(channel)
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
		c.dispatch(msg)
		return nil
	}
}
