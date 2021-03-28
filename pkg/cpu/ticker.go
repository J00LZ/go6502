package cpu

import (
	"time"
)

type Ticker struct {
	T       chan time.Time
	TTicker *time.Ticker
}

func NewTicker(ticker *time.Ticker) *Ticker {
	t := Ticker{T: make(chan time.Time, 1), TTicker: ticker}
	go func() {
		for {
			zz := <-t.TTicker.C
			t.T <- zz
		}
	}()
	return &t
}
