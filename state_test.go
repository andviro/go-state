package state_test

import (
	"github.com/andviro/go-state"
	"golang.org/x/net/context"
	"testing"
)

type IntState int

func (s *IntState) One(ctx context.Context) state.Func {
	*s += 1
	return s.Two
}

func (s *IntState) Two(ctx context.Context) state.Func {
	*s += 10
	return s.Three
}

func (s *IntState) Three(ctx context.Context) state.Func {
	*s += 100
	return nil
}

func TestRun(t *testing.T) {
	s := IntState(0)

	err := state.Run(context.TODO(), s.One)
	if err != nil {
		t.Error(err)
	}
	if s != 111 {
		t.Errorf("Invalid final value:", s)
	}
}

func TestHook(t *testing.T) {
	s := IntState(0)
	var temp string

	hook := func(c context.Context) error {
		temp += "-" + state.Name(c)
		return nil
	}

	err := state.Run(context.TODO(), s.One, hook)
	if err != nil {
		t.Error(err)
	}
	if temp != "-One-Two-Three" {
		t.Error("Invalid final value:", temp)
	}
}

func TestAllHooksAreRun(t *testing.T) {
	s := IntState(0)
	var temp string

	hook1 := func(c context.Context) error {
		temp += "1>" + state.Name(c) + ">"
		return nil
	}
	hook2 := func(c context.Context) error {
		temp += "2>" + state.Name(c) + ">"
		return nil
	}

	err := state.Run(context.TODO(), s.One, hook1, hook2)
	if err != nil {
		t.Error(err)
	}
	if temp != "1>One>2>One>1>Two>2>Two>1>Three>2>Three>" {
		t.Error("Invalid final value:", temp)
	}
}

func TestPanicRecovery(t *testing.T) {
	s := IntState(0)
	var temp string

	hook1 := func(c context.Context) error {
		temp += "1>" + state.Name(c) + ">"
		if s > 10 {
			panic("Aaaargh!!!")
		}
		return nil
	}
	hook2 := func(c context.Context) error {
		temp += "2>" + state.Name(c) + ">"
		return nil
	}

	err := state.Run(context.TODO(), s.One, hook1, hook2)
	if err.Error() != "Aaaargh!!!" {
		t.Errorf("Panic not recovered correctly: %v", err)
	}

	if _, ok := err.(state.Error); !ok {
		t.Errorf("Invalid type for error returned")
	}

	if temp != "1>One>2>One>1>Two>2>Two>1>Three>" {
		t.Error("Invalid final value:", temp)
	}
}
