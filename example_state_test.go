package state_test

import (
	"fmt"
	"github.com/andviro/go-state"
)

// Some user type holding state context
type Counter struct {
	X        int
	Messages chan bool
	Result   chan int
}

// Here go some state functions that comprise the machine

func (c *Counter) Start() state.Func {
	c.X = 0
	return c.Main // Start -> Main
}

func (c *Counter) Main() state.Func {
	for msg := range c.Messages {
		if msg {
			return c.Inc // Main -> Inc
		} else {
			return c.Dec // Main -> Dec
		}
	}
	return c.Stop // Main -> Stop
}

func (c *Counter) Inc() state.Func {
	c.X += 1
	return c.Main // Inc -> Main
}

func (c *Counter) Dec() state.Func {
	c.X -= 1
	return c.Main // Dec -> Main
}

func (c *Counter) Stop() state.Func {
	fmt.Println("Finished counting")
	c.Result <- c.X
	return nil // Stop -> END
}

func Example() {
	ctr := Counter{0, make(chan bool), make(chan int)}

	// Here we use hook function to do something on state change
	go state.Run(ctr.Start, func(s state.Func) error {
		fmt.Printf("Entered state %s, counter value is %d\n", s.Name(), ctr.X)
		return nil
	})

	// Send some messages to machine
	for _, msg := range []bool{true, false, false, true, true, false} {
		ctr.Messages <- msg
	}

	// Signal machine termination and wait for the result
	close(ctr.Messages)
	fmt.Printf("Final value is %d\n", <-ctr.Result)

	// Output:
	// Entered state Start, counter value is 0
	// Entered state Main, counter value is 0
	// Entered state Inc, counter value is 0
	// Entered state Main, counter value is 1
	// Entered state Dec, counter value is 1
	// Entered state Main, counter value is 0
	// Entered state Dec, counter value is 0
	// Entered state Main, counter value is -1
	// Entered state Inc, counter value is -1
	// Entered state Main, counter value is 0
	// Entered state Inc, counter value is 0
	// Entered state Main, counter value is 1
	// Entered state Dec, counter value is 1
	// Entered state Main, counter value is 0
	// Entered state Stop, counter value is 0
	// Finished counting
	// Final value is 0
}
