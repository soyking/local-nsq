package lnsq

import (
	"sync"
)

type Channel struct {
	Name string
	Stat *CountStat

	topicsLock sync.RWMutex
	topics     map[string]*Topic
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name:   name,
		topics: make(map[string]*Topic),
		Stat:   NewCountStat(),
	}
}

func (c *Channel) Subscribe(topic string, callback Callback, maxInFlight int, concurrency int) *Topic {
	c.topicsLock.Lock()
	defer c.topicsLock.Unlock()

	t, ok := c.topics[topic]
	if !ok {
		t = NewTopic(topic, maxInFlight)
		c.topics[topic] = t
	}

	t.Subscribe(callback, concurrency)
	return t
}

func (c *Channel) Dispatch(msg interface{}) {
	c.topicsLock.RLock()
	defer c.topicsLock.RUnlock()

	c.Stat.IncrCount()
	for _, topic := range c.topics {
		topic.Dispatch(msg)
	}
}

func (c *Channel) TopicsStats() map[string]int64 {
	c.topicsLock.RLock()
	defer c.topicsLock.RUnlock()

	topicsStats := make(map[string]int64)
	for name, topic := range c.topics {
		topicsStats[name] = topic.Stat.Count()
	}
	return topicsStats
}

func (c *Channel) CallbacksStats(topic string) map[string]int64 {
	c.topicsLock.RLock()
	defer c.topicsLock.RUnlock()

	if t, ok := c.topics[topic]; !ok {
		return map[string]int64{}
	} else {
		return t.CallbacksStats()
	}
}
