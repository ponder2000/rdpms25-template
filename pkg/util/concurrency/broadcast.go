package concurrency

import (
	"context"
	"time"
)

func Broadcast[T any](ctx context.Context, msg T, channels []chan T, timeouts []time.Duration, optionalTimeoutHandler func(index int, msg T)) {
	for i := range channels {
		go func(channel chan T, index int) {
			t := time.NewTimer(timeouts[index])
			select {
			case channel <- msg:
				if !t.Stop() {
					<-t.C
				}
			case <-ctx.Done():
				if !t.Stop() {
					<-t.C
				}
			case <-t.C:
				if optionalTimeoutHandler != nil {
					optionalTimeoutHandler(index, msg)
				}
			}
		}(channels[i], i)
	}
}
