package fsm

import (
	"errors"
	"fmt"
)

var (
	// ErrFsmLocked state machine is locked
	ErrFsmLocked = errors.New("state machine is locked")
	// ErrFsmUnLocked state machine is unlocked
	ErrFsmUnLocked = errors.New("state machine is unlocked")
	// ErrFsmBusy fsm is busy
	ErrFsmBusy = errors.New("State machine is busy")
	// ErrFsmEnded fsm has ended
	ErrFsmEnded = errors.New("State Machine has ended")
)

// ErrUnKnownState unknown state
func ErrUnKnownState(state State) error {
	return fmt.Errorf("Unknwon State %v", state)
}

// ErrUnKnownEvent unknown event
func ErrUnKnownEvent(event Event) error {
	return fmt.Errorf("Unknwon event %v", event)
}

// ErrDuplicateEvent duplicate event
func ErrDuplicateEvent(event Event) error {
	return fmt.Errorf("duplicate event %v", event)
}

// ErrDuplicateState duplicate state
func ErrDuplicateState(state State) error {
	return fmt.Errorf("duplicate State %v", state)
}

// ErrDuplicateTransition unknown transition
func ErrDuplicateTransition(from, to State) error {
	return fmt.Errorf("duplicate Transistion from %v to %v", from, to)
}

// ErrUnKnownTransition unknown transition
func ErrUnKnownTransition(from, to State) error {
	return fmt.Errorf("Unsupported Transistion from %v to %v", from, to)
}
