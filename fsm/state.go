package fsm

// State state for state machine
type State interface {
	String() string
}

// StateProvider interface which provides state information
type StateProvider interface {
	CurrentState() State
	SetState(state State)
}

// StateExists check if state exists in fsm
func (fsm *StateMachine) StateExists(state State) bool {
	_, ok := fsm.index[state]
	return ok
}

// IsEndState check if state is end state
func (fsm *StateMachine) IsEndState(state State) bool {
	return !fsm.StateTransitionsExists(state)
}

// AddState add state in state machine
func (fsm *StateMachine) AddState(states ...State) error {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	if fsm.lock {
		return ErrFsmLocked
	}
	for _, state := range states {
		if _, ok := fsm.index[state]; ok {
			return ErrDuplicateState(state)
		}
		fsm.index[state] = len(fsm.states)
		fsm.states = append(fsm.states, state)
	}
	return nil
}
