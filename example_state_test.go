package state_test

import (
	"context"
	"fmt"
	"gopkg.in/andviro/go-state.v2"
)

// Counter is a user type holding integer variable.
// Each instance of the Counter is a separate state machine.
type Counter int

// Context provides an environment for state machine.
// Different machines can use common initial context, to share a message channel, for example.

// InPort extracts message input port from state context
func InPort(ctx context.Context) chan bool {
	return ctx.Value("InPort").(chan bool)
}

// OutPort extracts result output port from state context
func OutPort(ctx context.Context) chan Counter {
	return ctx.Value("OutPort").(chan Counter)
}

// Here go some state functions that comprise the machine

func (c *Counter) Start(ctx context.Context) state.Func {
	*c = 0
	return c.Main // Start -> Main
}

func (c *Counter) Main(ctx context.Context) state.Func {
	for {
		select {
		case msg := <-InPort(ctx):
			if msg {
				return c.Inc // Main -> Inc
			} else {
				return c.Dec // Main -> Dec
			}
		case <-ctx.Done():
			return c.Stop // Main -> Stop
		}
	}
}

func (c *Counter) Inc(ctx context.Context) state.Func {
	*c += 1
	return c.Main // Inc -> Main
}

func (c *Counter) Dec(ctx context.Context) state.Func {
	*c -= 1
	return c.Main // Dec -> Main
}

func (c *Counter) Stop(ctx context.Context) state.Func {
	fmt.Println("Finished counting")
	OutPort(ctx) <- *c
	return nil // Stop -> END
}

func Example() {
	ctr := Counter(-1)
	initialCtx, cancel := context.WithCancel(context.TODO())
	ctx := context.WithValue(
		context.WithValue(initialCtx, "InPort", make(chan bool)),
		"OutPort",
		make(chan Counter),
	)

	// Here we use hook function to do something on state change
	go state.Run(ctx, ctr.Start, func(c context.Context) error {
		switch n := state.Name(c); n {
		case "Start", "Stop", "Main":
			fmt.Printf("%s: %d\n", n, ctr)
		case "Inc":
			fmt.Println("+1")
		case "Dec":
			fmt.Println("-1")
		}
		return nil
	})

	// Send some messages to machine
	for _, msg := range []bool{true, false, false, true, true, false} {
		InPort(ctx) <- msg
	}

	// Signal machine termination and wait for the result
	cancel()
	fmt.Printf("Final value is %d\n", <-OutPort(ctx))

	// Output:
	// Start: -1
	// Main: 0
	// +1
	// Main: 1
	// -1
	// Main: 0
	// -1
	// Main: -1
	// +1
	// Main: 0
	// +1
	// Main: 1
	// -1
	// Main: 0
	// Stop: 0
	// Finished counting
	// Final value is 0
}
