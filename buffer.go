package main

import (
	"context"
	"fmt"
	"time"
)

func buffer[T any](ctx context.Context, stream <-chan T, maxItems int, maxWait time.Duration) <-chan []T {
	batch := make(chan []T, 1)
	buffered := make([]T, 0, maxItems)

	go func() {
		timer := time.NewTimer(maxWait)
		defer timer.Stop()
		defer close(batch)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("Timed out")
				return
			case v := <-stream:
				buffered = append(buffered, v)
				if len(buffered) >= maxItems {
					var bufCopy []T
					bufCopy = append(bufCopy, buffered...)
					batch <- bufCopy
					buffered = buffered[:0]
					timer.Reset(maxWait)
				}
			case <-timer.C:
				var bufCopy []T
				bufCopy = append(bufCopy, buffered...)
				batch <- bufCopy
				buffered = buffered[:0]
				timer.Reset(maxWait)
			}
		}
	}()
	return batch

}
