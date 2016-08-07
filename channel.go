package lnsq

import "sync"

type Channel struct {
	Name       string
	topicsLock sync.RWMutex
	topics     map[string]*Topic
	Stats      *CountStats
}

func NewChannel(name string, stats *CountStats) *Channel {
	return &Channel{
		Name:   name,
		topics: make(map[string]*Topic),
		Stats:  stats,
	}
}

func (c *Channel) Subscribe(topic string, callback Callback, maxInFlight int, concurrency int) *Topic {
	c.topicsLock.Lock()
	defer c.topicsLock.Unlock()

	t, ok := c.topics[topic]
	if !ok {
		t = NewTopic(topic, maxInFlight, c.Stats.NewSubStats(topic))
		c.topics[topic] = t
	}

	t.Subscribe(callback, concurrency)
	return t
}

func (c *Channel) Dispatch(msg interface{}) {
	c.topicsLock.RLock()
	defer c.topicsLock.RUnlock()

	c.Stats.IncrCount()
	for _, topic := range c.topics {
		topic.Dispatch(msg)
	}
}

func (c *Channel) TopicsStats() map[string]int64 {
	c.Stats.RLock()
	defer c.Stats.RUnlock()

	topicsStats := make(map[string]int64)
	for topic, stats := range c.Stats.SubStats {
		topicsStats[topic] = stats.Count
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
