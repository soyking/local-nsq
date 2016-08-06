package lnsq

type Callback func(interface{})

type topic struct {
	name string
	ch   chan interface{}
}

func newTopic(name string, maxInFlight int) *topic {
	return &topic{
		name: name,
		ch:   make(chan interface{}, maxInFlight),
	}
}

func (t *topic) subscribe(c Callback, concurrency int) {
	for i := 0; i < concurrency; i++ {
		go t.callback(c)
	}
}

func (t *topic) callback(c Callback) {
	for msg := range t.ch {
		c(msg)
	}
}

func (t *topic) dispatch(msg interface{}) {
	t.ch <- msg
}
