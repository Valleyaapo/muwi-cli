package main

import "context"

// randNumFetcher := func() int { return rand.Intn(500000000) }
// randIntStream := repeatFunc(ctx, randNumFetcher)

func primeFinder(ctx context.Context, randIntStram <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}
	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-ctx.Done():
				return
			case randomInt := <-randIntStram:
				if isPrime(randomInt) {
					primes <- randomInt
				}
			}
		}
	}()
	return primes
}

func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
			}

		}
	}()
	return taken
}
func repeatFunc[T any](ctx context.Context, fn func() T) <-chan T {
	stream := make(chan T, 5)
	go func() {
		defer close(stream)
		for {
			select {
			case <-ctx.Done():
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}
