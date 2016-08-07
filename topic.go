package lnsq

type Callback func(interface{})

type Topic struct {
	Name  string
	ch    chan interface{}
	Stats *CountStats
}

func NewTopic(name string, maxInFlight int, stats *CountStats) *Topic {
	return &Topic{
		Name:  name,
		ch:    make(chan interface{}, maxInFlight),
		Stats: stats,
	}
}

func (t *Topic) Subscribe(c Callback, concurrency int) {
	for i := 0; i < concurrency; i++ {
		go t.Callback(c)
	}
}

func (t *Topic) Callback(c Callback) {
	callbackStats := t.Stats.NewSubStats("")
	for msg := range t.ch {
		c(msg)
		callbackStats.IncrCount()
	}
}

func (t *Topic) Dispatch(msg interface{}) {
	t.Stats.IncrCount()
	t.ch <- msg
}

func (t *Topic) CallbacksStats() map[string]int64 {
	callbacksStats := make(map[string]int64)
	for topic, stats := range t.Stats.SubStats {
		callbacksStats[topic] = stats.Count
	}
	return callbacksStats
}
