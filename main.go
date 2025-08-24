package main

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"
)

func main() {
	now := time.Now()
	fmt.Println(("Starting muwi application..."))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	orderChan := readCSV("online_furniture_retailer.csv", ctx)

	CPUcount := runtime.NumCPU()
	orderEditorChannels := make([]<-chan Order, CPUcount)

	for i := range CPUcount {
		orderEditorChannels[i] = orderEditor(ctx, orderChan)
	}
	f := fanIn(ctx, orderEditorChannels...)
	batches := buffer(ctx, f, 2, 500*time.Millisecond)
	items := takeItems(ctx, batches, 10)
	for v := range items {
		fmt.Println("One order: ", v)
	}
	fmt.Printf("Took %s", time.Since(now))
}

func orderEditor(ctx context.Context, orders <-chan Order) <-chan Order {
	editedOrders := make(chan Order)
	cutFloat := func(f float64) float64 {
		return math.Round(f*100) / 100
	}
	go func() {
		defer close(editedOrders)
		for {
			select {
			case <-ctx.Done():
				return
			case order, ok := <-orders:
				if !ok {
					return
				}
				order.CustomerRating = cutFloat(order.CustomerRating)
				order.TotalAmount = cutFloat(order.TotalAmount)
				editedOrders <- order
			}
		}
	}()
	return editedOrders
}

func fanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	fStream := make(chan T)

	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				return
			case fStream <- i:
			}
		}
	}
	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}
	go func() {
		wg.Wait()
		close(fStream)
	}()
	return fStream
}

// Flatten batches ([]T) into items (T) and stop after n items.
func takeItems[T any](ctx context.Context, batches <-chan []T, n int) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		remaining := n
		// optional local carry for overshoot; we only need it within this call
		var carry []T

		emit := func(v T) bool { // returns true if we should stop
			select {
			case <-ctx.Done():
				return true
			case out <- v:
				remaining--
				return remaining == 0
			}
		}

		// drain carry first if present
		for len(carry) > 0 && remaining > 0 {
			v := carry[0]
			carry = carry[1:]
			if emit(v) {
				return
			}
		}

		for remaining > 0 {
			select {
			case <-ctx.Done():
				return
			case batch, ok := <-batches:
				if !ok {
					return // upstream finished; we’re done with whatever we emitted
				}
				for i := 0; i < len(batch) && remaining > 0; i++ {
					if emit(batch[i]) {
						return
					}
				}
				// if batch had more than we needed, stash the rest locally (optional)
				if len(batch) > 0 && len(batch) > n {
					// Not strictly needed unless you plan to reuse leftovers later.
					// For a one-shot, you can ignore carry entirely.
				}
			}
		}
	}()
	return out
}
