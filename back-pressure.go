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
	quite               chan struct{}
}

func NewBackPressure[T any](maxSizeChannel, thresholdPercentage int64, wait time.Duration) *BackPressure[T] {
	log.Println("making the backpressure")
	b := &BackPressure[T]{
		maxSizeChannel:      maxSizeChannel,
		thresholdPercentage: thresholdPercentage,
		channel:             make(chan T, maxSizeChannel),
		wait:                wait,
		quite:               make(chan struct{}),
	}
	go b.Wait()
	go b.ShotDown()
	return b
}
func (b *BackPressure[T]) ShotDown() {

	<-b.quite

	if len(b.channel) == 0 {
		log.Println("BACK PRESSURE: CHANNEL IS EMPTY NOW I CLOSE IT")
		close(b.channel)
	}

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

func (b *BackPressure[T]) Close() {
	close(b.quite)
}
