package lnsq

type Callback func(interface{})

type topic struct {
	name  string
	ch    chan interface{}
	stats *CountStats
}

func newTopic(name string, maxInFlight int, stats *CountStats) *topic {
	return &topic{
		name:  name,
		ch:    make(chan interface{}, maxInFlight),
		stats: stats,
	}
}

func (t *topic) subscribe(c Callback, concurrency int) {
	for i := 0; i < concurrency; i++ {
		go t.callback(c)
	}
}

func (t *topic) callback(c Callback) {
	subStats := t.stats.NewSubStats("")
	for msg := range t.ch {
		c(msg)
		subStats.Count()
	}
}

func (t *topic) dispatch(msg interface{}) {
	t.ch <- msg
}
