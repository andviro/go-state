package state

import (
	"fmt"
	"golang.org/x/net/context"
	"reflect"
	"regexp"
	"runtime"
)

type stateKeyType int

const (
	stateKey stateKeyType = iota
)

// Func is a basic building block of state machine. It's a simple function that does some work,
// maybe listens to channels and returns next state, based on arbitrary conditions.
// The context is passed around as the state machine progresses.
// Each time the state is changed it is injected into context
type Func func(context.Context) Func

// Hook is a function that's run when machine is about to enter certain state.
// Context parameter allows to extract current state name using `Name()`
type Hook func(context.Context) error

var nameRe = regexp.MustCompile(`(\w+)([-][^-]*)?$`)

// Name returns name of current state extracted from context
func Name(ctx context.Context) string {
	f, ok := ctx.Value(stateKey).(Func)
	if !ok {
		return "<Undefined>"
	}
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	return nameRe.FindStringSubmatch(name)[1]
}

// Run starts the state machine with the provided context and initial Func.
// Each state produces next one until last Func returns nil.
// Each time the state is changed the transition hooks are run on current context.
// If transition hook returns non-nil error, state machine terminates and returns the error.
// All panics in states and hooks are recovered and converted to errors.
func Run(ctx context.Context, initial Func, hooks ...Hook) (err error) {
	defer func() {
		if e := recover(); e != nil {
			var ok bool
			err, ok = e.(error)
			if !ok {
				err = fmt.Errorf("Panic: %v", e)
			}
		}
	}()

	for initial != nil {
		nextC := context.WithValue(ctx, stateKey, initial)
		for _, h := range hooks {
			err = h(nextC)
			if err != nil {
				return
			}
		}
		initial = initial(nextC)
	}
	return
}