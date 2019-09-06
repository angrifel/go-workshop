package main

import (
	"fmt"
	"sync"
)

const goroutines = 10
const countLimit = 10000

func main() {
	counter := 0
	wg := sync.WaitGroup{}
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for j := 0; j < countLimit; j++ {
				counter++
			}
		}()
	}

	wg.Wait()

	fmt.Printf("counter: %d\n", counter)
}
