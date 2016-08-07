package lnsq

import (
	"math/rand"
	"time"
)

var (
	randBase    = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	randBaseLen = len(randBase)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randStatID() string {
	return randStr(10)
}

func randStr(n int) string {
	r := make([]byte, n)
	for i := 0; i < n; i++ {
		r[i] = randBase[rand.Intn(randBaseLen)]
	}
	return string(r)
}
