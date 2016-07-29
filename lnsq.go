package lnsq

import (
	"sync"
)

type Callback func(interface{}) error

type LocalNSQ struct {
	sync.Mutex
	eventCh   chan *Event
	callbacks map[string][]Callback
}

type Event struct {
	Topic   string
	Content interface{}
}

const (
	defaultSize = 10
)

func NewLocalNSQ(size ...int) *LocalNSQ {
	s := defaultSize
	if len(size) != 0 {
		s = size[0]
	}

	l := &LocalNSQ{
		eventCh:   make(chan *Event, s),
		callbacks: make(map[string][]Callback),
	}
	go l.receive()
	return l
}

func (l *LocalNSQ) receive() {
	for e := range l.eventCh {
		topic := e.Topic
		go dispatchTopic(l.callbacks[topic], e.Content)
	}
}

func dispatchTopic(callbacks []Callback, content interface{}) {
	for i := range callbacks {
		callbacks[i](content)
	}
}

func (l *LocalNSQ) Dispatch(topic string, content interface{}) {
	l.eventCh <- &Event{Topic: topic, Content: content}
}

func (l *LocalNSQ) Subscribe(topic string, c Callback) {
	l.Lock()
	defer l.Unlock()
	if l.callbacks[topic] == nil {
		l.callbacks[topic] = []Callback{}
	}
	l.callbacks[topic] = append(l.callbacks[topic], c)
}

func (l *LocalNSQ) Close() {
	close(l.eventCh)
}
