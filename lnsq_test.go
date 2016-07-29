package lnsq

import (
	"testing"
	"time"
)

func TestLocalNSQ(t *testing.T) {
	l := NewLocalNSQ()
	defer l.Close()
	l.Subscribe("test", printCallback)
	l.Dispatch("test", 1)
	time.Sleep(1 * time.Second)
}

func printCallback(v interface{}) error {
	println(v.(int))
	return nil
}