package fsm

// Action action for state machine
type Action func(obj StateProvider, from, to State) error

// Transition transition from one state to another
type Transition map[State]Action

// StateTransitionsExists check if transition exists in fsm
func (fsm *StateMachine) StateTransitionsExists(state State) bool {
	_, ok := fsm.transitions[state]
	return ok
}

// TransitionExists check is transition exists
func (fsm *StateMachine) TransitionExists(from, to State) bool {
	_, ok := fsm.transitions[from][to]
	return ok
}

// TransitionAction get transition action
func (fsm *StateMachine) TransitionAction(from, to State) (Action, bool) {
	v, ok := fsm.transitions[from][to]
	return v, ok
}

// AddTransition add transition
func (fsm *StateMachine) AddTransition(from, to State, action Action) error {
	fsm.mutex.Lock()
	defer fsm.mutex.Unlock()
	if fsm.lock {
		return ErrFsmLocked
	}
	if !fsm.StateExists(from) {
		return ErrUnKnownState(from)
	}
	if !fsm.StateExists(to) {
		return ErrUnKnownState(to)
	}
	if _, ok := fsm.transitions[from][to]; ok {
		return ErrDuplicateTransition(from, to)
	}
	if !fsm.StateTransitionsExists(from) {
		fsm.transitions[from] = map[State]Action{}
	}
	fsm.transitions[from][to] = action
	return nil
}
