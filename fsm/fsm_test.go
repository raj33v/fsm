package fsm_test

import (
	"fmt"
	"testing"

	"bitbucket.org/raj33v/srotsa-media-api/app/fsm"
)

type myState int

const (
	state1 myState = iota + 1
	state2
	state3
	state4
	state5
	state6
	stateLast
)

func (s myState) String() string {
	if s >= stateLast {
		return "stateLast"
	}
	return fmt.Sprintf("state%d", s)
}

type testObject struct {
	state myState
}

func (t *testObject) CurrentState() fsm.State {
	return t.state
}
func (t *testObject) SetState(s fsm.State) {
	t.state = s.(myState)
}

func TestFsm(t *testing.T) {
	machine := fsm.New(state1)
	machine.AddState(state2, state3, state4, state5, state6, stateLast)
	fn := func(obj fsm.StateProvider, from, to fsm.State) error {
		fmt.Printf("Transition from %v to %v\n", from, to)
		return nil
	}
	fn1 := func(obj fsm.StateProvider, state fsm.State) {
		fmt.Printf("State is %v\n", state)
	}
	go func() {
		ch := machine.OpenChannel()
		for {
			state := <-ch
			if state == nil {
				fmt.Println("ended")
				return
			}
			fmt.Printf("From Async Event State %v acheieved\n", state)
		}
	}()
	machine.AddTransition(state1, state2, fn)
	machine.AddTransition(state2, state3, fn)
	machine.AddTransition(state2, state4, fn)
	machine.AddTransition(state3, state4, fn)
	machine.AddTransition(state4, state5, fn)
	machine.AddTransition(state5, state6, fn)
	machine.AddTransition(state6, stateLast, nil)
	machine.RegisterEvent(fsm.EventState, state5, fn1)

	obj := &testObject{state1}
	machine.Lock()
	if err := machine.Run(obj, state2); err != nil {
		t.Fail()
		t.Errorf("%v", err)
	}
	if err := machine.Run(obj, state3); err != nil {
		t.Fail()
		t.Errorf("%v", err)
	}
	if err := machine.Run(obj, state4); err != nil {
		t.Fail()
		t.Errorf("%v", err)
	}
	if err := machine.Run(obj, state5); err != nil {
		t.Fail()
		t.Errorf("%v", err)
	}
	if err := machine.Run(obj, state6); err != nil {
		t.Fail()
		t.Errorf("%v", err)
	}
	if err := machine.Run(obj, stateLast); err != nil {
		t.Fail()
		t.Errorf("%v", err)
	}
}
