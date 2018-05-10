package fsm

// Event event
type Event int

// EventHandler event handler
type EventHandler func(StateProvider, State)

// EventChannel channel for async event notifications
type EventChannel chan *EventInfo

const (
	// EventBegin fsm has started
	EventBegin Event = iota
	// EventEnded fsm has ended
	EventEnded
	// EventStateEntering state machine is entering a state
	EventStateEntering
	// EventState state machine has entered a state
	EventState
	//EventStateLeaving state machine is leaving a state
	EventStateLeaving
	// EventStateLeft state machine has left the state
	EventStateLeft
	// EventInvalid invalid state event
	EventInvalid
)

const (
	// EventCount count of events
	EventCount int = int(EventInvalid)
)

// EventInfo event info for asynchronous events
type EventInfo struct {
	Event  Event
	State  State
	Source StateProvider
}

// EventManager event manager
type EventManager struct {
	handlers map[State]*[EventCount]EventHandler
	channel  EventChannel
}

// EventManager get fsm event manager
func (fsm *StateMachine) EventManager() *EventManager {
	return &EventManager{
		handlers: map[State]*[EventCount]EventHandler{},
	}
}

// RegisterEvent register event
func (e *EventManager) RegisterEvent(event Event, state State, handler EventHandler) error {
	if int(event) >= EventCount {
		return ErrUnKnownEvent(event)
	}
	if e.handlers[state] == nil {
		e.handlers[state] = &[EventCount]EventHandler{}
	}
	if e.handlers[state][event] != nil {
		return ErrDuplicateEvent(event)
	}
	e.handlers[state][event] = handler
	return nil
}

// Open return channel
func (e *EventManager) Open() EventChannel {
	if e.channel != nil {
		return e.channel
	}
	ch := make(EventChannel)
	e.channel = ch
	return ch
}

// Close close channel
func (e *EventManager) Close() {
	if e.channel != nil {
		ch := e.channel
		e.channel = nil
		close(ch)
	}
}

// Raise raise an event
func (e *EventManager) Raise(obj StateProvider, event Event, state State) {
	if int(event) >= EventCount {
		return
	}
	if e.handlers[state] != nil && e.handlers[state][event] != nil {
		e.handlers[state][event](obj, state)
	}
	if e.channel != nil {
		if event != EventEnded {
			info := &EventInfo{Event: event, Source: obj, State: state}
			e.channel <- info
		} else {
			e.channel <- nil
		}
	}
}

// EnteringState entering a state
func (e *EventManager) EnteringState(obj StateProvider, state State) {
	e.Raise(obj, EventStateEntering, state)
}

// State state achieved
func (e *EventManager) State(obj StateProvider, state State) {
	e.Raise(obj, EventState, state)
}

// LeavingState leaving state
func (e *EventManager) LeavingState(obj StateProvider, state State) {
	e.Raise(obj, EventStateLeaving, state)
}

// LeftState left state
func (e *EventManager) LeftState(obj StateProvider, state State) {
	e.Raise(obj, EventStateLeft, state)
}

// Begin fsm has started with state
func (e *EventManager) Begin(obj StateProvider, state State) {
	e.Raise(obj, EventBegin, state)
}

// Ended fsm has ended with state
func (e *EventManager) Ended(obj StateProvider, state State) {
	e.Raise(obj, EventEnded, state)
}
