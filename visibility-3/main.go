package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	_initialValue = int64(1)
	_increment    = int64(2)
	_lastValue    = int64(20000)
)

func main() {

	variableSet := int64(0)
	a := int64(-1)
	b := int64(0)
	done := int64(0)

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		_nextValue := _increment
		for _nextValue <= _lastValue {
			mutex.Lock()
			// read shared variables into local variables
			_a := atomic.LoadInt64(&a)
			_b := atomic.LoadInt64(&b)

			// evaluate local variables
			if _b == _nextValue {
				if _b == _a+1 {
					_nextValue = _b + _increment
					atomic.StoreInt64(&variableSet, 0) // we shared variable to signal the other goroutine it can continue with next value
				} else {
					fmt.Printf("memory visibility issue with a= %d, b= %d\n", _a, _b)
					atomic.StoreInt64(&done, 1)
				}
			}

			mutex.Unlock()
		}
		atomic.StoreInt64(&done, 1)

		fmt.Printf("comparer out\n")
	}()

	// time.Sleep(1 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_counter := int64(-1) // local variable

		for {
			mutex.Lock()
			_done := atomic.LoadInt64(&done)
			if _done != 0 {
				break
			}

			// read global variable
			_variableSet := variableSet
			if _variableSet == 0 {
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

			mutex.Unlock()
		}

		fmt.Printf("incrementer out\n")
	}()

	wg.Wait()
}
