local-nsq
=========
[NSQ](http://nsq.io/) Local and Simple Version

## Usage

```
go get github.com/soyking/local-nsq
```

## Definition

`channel`, `topic` same as NSQ

`maxInFlight` the length of the chan's buffer of a topic
 
`concurrency` the number of cunsummer(callback) working at one time

## Example

```go
package main

import (
	"github.com/soyking/local-nsq"
	"sync"
)

var wg sync.WaitGroup

func intCallback(v interface{}) {
	println(v.(int))
	wg.Done()
}

func main() {
	l := lnsq.NewLocalNSQ()
	wg.Add(1)
	l.Subscribe("channel", "topic", intCallback, 1, 1)
	l.Dispatch("channel", 17)
	wg.Wait()
}

```

