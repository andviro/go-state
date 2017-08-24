package state_test

import (
	"context"
	"github.com/andviro/go-state"
	"testing"
)

var chainResults string

func Mw1(next state.Func) state.Func {
	return func(ctx context.Context) state.Func {
		chainResults += "1>"
		return next
	}
}

func Mw2(next state.Func) state.Func {
	return func(ctx context.Context) state.Func {
		chainResults += "2>"
		return next
	}
}

func Mw3(next state.Func) state.Func {
	return func(ctx context.Context) state.Func {
		chainResults += "3>"
		return next
	}
}

func TestChain(t *testing.T) {
	s := IntState(0)
	mw := state.Chain(Mw1, Mw2, Mw3)
	chainResults = ""

	state.Run(context.TODO(), mw(s.One))
	if chainResults != "1>2>3>" {
		t.Errorf("Invalid final value: %s", chainResults)
	}
}

func TestUse(t *testing.T) {
	s := IntState(0)
	mw := state.Chain(Mw1, Mw2, Mw3)
	mw = mw.Use(Mw1)
	chainResults = ""

	state.Run(context.TODO(), mw(s.One))
	if chainResults != "1>2>3>1>" {
		t.Errorf("Invalid final value: %s", chainResults)
	}
}
