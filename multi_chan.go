package lnsq

type MultiChan struct {
	cIndex int
	c      []chan interface{}
}

func NewMultiChan(n int, size int) *MultiChan {
	c := []chan interface{}{}
	for i := 0; i < n; i++ {
		c = append(c, make(chan interface{}, size))
	}
	return &MultiChan{cIndex: -1, c: c}
}

func (c *MultiChan) Dispatch(event interface{}) {

}
