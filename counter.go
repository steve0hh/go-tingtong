package main

import "time"

type Counter struct {
	count chan int
	reset chan bool
	done  chan bool
}

func newCounter() *Counter {
	return &Counter{
		count: make(chan int),
		reset: make(chan bool),
		done:  make(chan bool),
	}
}

func (c *Counter) start(d time.Duration) {
	count := 0
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			count++
			c.count <- count
		case <-c.reset:
			count = 0
		case <-c.done:
			ticker.Stop()
		}
	}
}
