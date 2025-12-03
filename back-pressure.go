package BackPressure

import (
	"log"
	"time"
)

type BackPressure[T any] struct {
	maxSizeChannel      int64
	thresholdPercentage int64
	channel             chan T
	wait                time.Duration
}

func NewBackPressure[T any](maxSizeChannel, thresholdPercentage int64, wait time.Duration) *BackPressure[T] {
	log.Println("making the backpressure")
	b := &BackPressure[T]{
		maxSizeChannel:      maxSizeChannel,
		thresholdPercentage: thresholdPercentage,
		channel:             make(chan T, maxSizeChannel),
		wait:                wait,
	}
	go b.Wait()
	return b
}

func (b *BackPressure[T]) Add(item T) {
	b.channel <- item
}
func (b *BackPressure[T]) Wait() {
	for {
		if int64(len(b.channel)) >= b.thresholdPercentage*b.maxSizeChannel/100 {
			log.Println("channel has reached it's threshold")
			time.Sleep(b.wait)
		} else {
			log.Println("channel is open")
		}
	}
}

func (b *BackPressure[T]) Out() chan T {
	return b.channel
}
