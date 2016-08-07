package lnsq

type Callback func(interface{})

type CallbackWrapper struct {
	callback Callback
	Stat     *CountStat
}

func NewCallbackWrapper(callback Callback) *CallbackWrapper {
	return &CallbackWrapper{
		callback: callback,
		Stat:     NewCountStat(),
	}
}

func (c *CallbackWrapper) Handler(msg interface{}) {
	c.callback(msg)
	c.Stat.IncrCount()
}
