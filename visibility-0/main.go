package main

import (
	"fmt"
	"sync"
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
	done := false

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		_nextValue := _increment
		for _nextValue <= _lastValue {
			// read shared variables into local variables
			_a := a
			_b := b

			// evaluate local variables
			if _b == _nextValue {
				if _b == _a+1 {
					_nextValue = _b + _increment
					variableSet = 0 // we shared variable to signal the other goroutine it can continue with next value
				} else {
					fmt.Printf("memory visibility issue with a= %d, b= %d\n", _a, _b)
					done = true
				}
			}
		}
	}()

	// time.Sleep(1 * time.Second)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_counter := int64(-1) // local variable

		for !done {
			// read global variable
			_variableSet := variableSet
			if _variableSet == 0 {
				// local variable calculation
				_a := a
				_b := b
				_a += _increment
				_b += _increment
				_counter++

				// set shared variables here
				a = _a
				b = _b
				variableSet = 1
			}
		}
	}()

	wg.Wait()
}
