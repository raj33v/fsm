package fsm

import (
	"sync"
)

// StateMachine state machine
type StateMachine struct {
	start       State
	states      []State
	index       map[State]int
	transitions map[State]Transition
	mutex       *sync.Mutex
	lock        bool
	event       *EventManager
	busy        bool
}

// New create new fsm
func New(start State) *StateMachine {
	fsm := &StateMachine{
		start:       start,
		index:       map[State]int{},
		transitions: map[State]Transition{},
		mutex:       &sync.Mutex{},
		lock:        false,
	}
	fsm.index[start] = len(fsm.states)
	fsm.states = append(fsm.states, start)
	fsm.event = fsm.EventManager()
	return fsm
}

// Lock lock state machine
func (fsm *StateMachine) Lock() error {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	if fsm.lock {
		return ErrFsmLocked
	}
	if fsm.busy {
		return ErrFsmBusy
	}
	fsm.lock = true
	return nil
}

// Unlock unlock state machine
func (fsm *StateMachine) Unlock() error {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	if !fsm.lock {
		return ErrFsmUnLocked
	}
	if fsm.busy {
		return ErrFsmBusy
	}
	fsm.lock = false
	return nil
}

// Ready check if fsm is ready
func (fsm *StateMachine) Ready() bool {
	return !fsm.busy
}

// Run run state machine on source
func (fsm *StateMachine) Run(source StateProvider, state State) error {
	if !fsm.lock {
		return ErrFsmUnLocked
	}
	fsm.busy = true
	defer func() { fsm.busy = false }()
	// Get current state
	current := source.CurrentState()
	if !fsm.StateExists(current) {
		return ErrUnKnownState(current)
	}
	// Check if end state
	if fsm.IsEndState(current) {
		// If end state is achieved then we can not continue
		return ErrFsmEnded
	}

	action, ok := fsm.TransitionAction(current, state)

	// Check if transition exists
	if !ok {
		return ErrUnKnownTransition(current, state)
	}
	// Check if fsm has just started
	if current == fsm.start {
		fsm.event.Begin(source, current)
	}
	// FSM is leaving from state
	fsm.event.LeavingState(source, state)
	// FSM is entering state
	fsm.event.EnteringState(source, state)
	source.SetState(state)
	if action != nil {
		if err := action(source, current, state); err != nil {
			source.SetState(current)
			return err
		}
	}
	// FSM has left the state
	fsm.event.LeftState(source, state)
	// FSM has enetered the state
	fsm.event.State(source, state)

	// Check if fsm has ended
	if fsm.IsEndState(state) {
		fsm.event.Ended(source, state)
	}
	return nil
}

// RegisterEvent register event
func (fsm *StateMachine) RegisterEvent(event Event, state State, handler EventHandler) error {
	return fsm.event.RegisterEvent(event, state, handler)
}

// OpenChannel get channel for events
// When a channel is created all state events and fsm events are
// sent to that channel so the fsm is blocked until some one
// is available to read from that channel
func (fsm *StateMachine) OpenChannel() EventChannel {
	return fsm.event.Open()
}

// CloseChannel close channel for events
func (fsm *StateMachine) CloseChannel() {
	fsm.event.Close()
}
