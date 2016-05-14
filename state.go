package state

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
)

// Func is a basic building block of state machine. It's a simple function that does some work,
// maybe listens to channels and returns next state, based on arbitrary conditions
type Func func() Func

// Hook is a function that's run when machine is about to enter certain state.
type Hook func(state Func) error

var nameRe = regexp.MustCompile(`(\w+)([-][^-]*)?$`)

// Name returns name of state function as a string
func (state Func) Name() string {
	name := runtime.FuncForPC(reflect.ValueOf(state).Pointer()).Name()
	return nameRe.FindStringSubmatch(name)[1]
}

// Run starts the state machine. Each state produces next one until last Func returns nil.
// Each time the state is changed the hooks are run.
// If transition hook returns non-nil error, state machine terminates and returns the error.
// All panics in states and hooks are recovered and converted to errors.
func Run(initial Func, hooks ...Hook) (err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("Panic: %v", e)
			}
		}
	}()

	for ; initial != nil; initial = initial() {
		for _, h := range hooks {
			err = h(initial)
			if err != nil {
				return
			}
		}
	}
	return
}
