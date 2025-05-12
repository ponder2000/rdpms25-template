package concurrency

import (
	"context"
	"sync"
)

func FanInContext[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	out := make(chan T)
	wg := sync.WaitGroup{}

	for i := range channels {
		wg.Add(1)
		go func(c <-chan T) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case val, ok := <-c:
					if !ok {
						return
					}
					select {
					case out <- val:
					case <-ctx.Done():
						return
					}
				}
			}
		}(channels[i])
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
