package engine

import "fmt"

func (e *E) Defer(action func()) {
	if len(e.defers) < cap(e.defers) {
		ln := len(e.defers)
		e.defers = e.defers[:ln+1]
		for i := ln; i > 0; i-- {
			e.defers[i] = e.defers[i-1]
		}
	}
}

type DeferredPanic struct {
	value interface{}
}

func (p DeferredPanic) Error() string {
	return fmt.Sprintf("%v", p.value)
}

func (e *E) executeDefers() (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = DeferredPanic{p}
		}
	}()
	return nil
}

