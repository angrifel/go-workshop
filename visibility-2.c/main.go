package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	_initialValue = int64(1)
	_increment    = int64(2)
	_lastValue    = int64(2000)
)

func main() {

	variableSet := int64(0)
	a := int64(-1)
	b := int64(0)
	done := int64(0)
	correctReadings := 0
	incorrectReadings := 0

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		_nextValue := _increment
		for _nextValue <= _lastValue {
			// read shared variables into local variables
			_a := atomic.LoadInt64(&a)
			_b := atomic.LoadInt64(&b)

			// evaluate local variables
			if _b == _nextValue {
				if _b == _a+1 {
					_nextValue = _b + _increment
					correctReadings++
					atomic.StoreInt64(&variableSet, 0) // we shared variable to signal the other goroutine it can continue with next value
				} else {
					incorrectReadings++
				}
			}
		}

		atomic.StoreInt64(&done, 1)
		fmt.Printf("comparer out\n")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		_counter := int64(-1) // local variable
		_readingsCount := 0
		for {
			_done := atomic.LoadInt64(&done)
			if _done != 0 {
				break
			}

			// read global variable
			_variableSet := atomic.LoadInt64(&variableSet)
			_readingsCount++
			if _variableSet == 0 {
				_readingsCount = 0

				// local variable calculation
				_a := atomic.LoadInt64(&a)
				_b := atomic.LoadInt64(&b)
				_a += _increment
				_b += _increment
				_counter++

				// set shared variables here
				atomic.StoreInt64(&a, _a)
				atomic.StoreInt64(&b, _b)
				atomic.StoreInt64(&variableSet, 1)
			}

			if _readingsCount > 100000 {
				panic("tired of waiting for an update, aborting...")
			}
		}

		fmt.Printf("incrementer out\n")
	}()

	wg.Wait()
	fmt.Printf("correct reading = %d, incorrect readings = %d\n", correctReadings, incorrectReadings)

}
