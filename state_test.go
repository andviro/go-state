package state_test

import (
	"github.com/andviro/go-state"
	"testing"
)

type IntState int

func (s *IntState) One() state.Func {
	*s += 1
	return s.Two
}

func (s *IntState) Two() state.Func {
	*s += 10
	return s.Three
}

func (s *IntState) Three() state.Func {
	*s += 100
	return nil
}

func TestRun(t *testing.T) {
	s := IntState(0)

	err := state.Run(s.One)
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

	hook := func(st state.Func) error {
		temp += "-" + st.Name()
		return nil
	}

	err := state.Run(s.One, hook)
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

	hook1 := func(st state.Func) error {
		temp += "1>" + st.Name() + ">"
		return nil
	}
	hook2 := func(st state.Func) error {
		temp += "2>" + st.Name() + ">"
		return nil
	}

	err := state.Run(s.One, hook1, hook2)
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

	hook1 := func(st state.Func) error {
		temp += "1>" + st.Name() + ">"
		if s > 10 {
			panic("Aaaargh!!!")
		}
		return nil
	}
	hook2 := func(st state.Func) error {
		temp += "2>" + st.Name() + ">"
		return nil
	}

	err := state.Run(s.One, hook1, hook2)
	if err.Error() != "Panic: Aaaargh!!!" {
		t.Error("Panic not recovered correctly", err)
	}
	if temp != "1>One>2>One>1>Two>2>Two>1>Three>" {
		t.Error("Invalid final value:", temp)
	}
}
