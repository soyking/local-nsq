package lnsq

import "sync"

type channel struct {
	name       string
	topicsLock sync.RWMutex
	topics     map[string]*topic
}

func newChannel(name string) *channel {
	return &channel{
		name:   name,
		topics: make(map[string]*topic),
	}
}

func (c *channel) subscribe(topic string, callback Callback, maxInFlight int, concurrency int) *topic {
	c.topicsLock.Lock()
	defer c.topicsLock.Unlock()

	t, ok := c.topics[topic]
	if !ok {
		t = newTopic(topic, maxInFlight)
		c.topics[topic] = t
	}

	t.subscribe(callback, concurrency)
	return t
}

func (c *channel) dispatch(msg interface{}) {
	c.topicsLock.RLock()
	defer c.topicsLock.RUnlock()

	for _, topic := range c.topics {
		topic.dispatch(msg)
	}
}
