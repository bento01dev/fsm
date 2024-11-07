package fsm

import (
	"errors"
	"sync"

	"golang.org/x/exp/slices"
)

var (
	ErrUnknownState      = errors.New("state not in list")
	ErrUnknownEvent      = errors.New("event not in list")
	ErrUnknownStart      = errors.New("start not in the state list")
	ErrTerminalState     = errors.New("fsm in terminal state.")
	ErrUnknownTransition = errors.New("transition with given event/state not in list")
)

type State int

type Event int

type EventState struct {
	Event Event
	Next  State
}

type Transition struct {
	State       State
	EventStates []EventState
}

type FSM struct {
	mu          sync.RWMutex
	current     State
	states      []State
	events      []Event
	transitions map[State]map[Event]State
}

func NewFSM(states []State, events []Event, start State, transitions ...Transition) (*FSM, error) {
	if !slices.Contains(states, start) {
		return nil, ErrUnknownStart
	}

	transitionMap := make(map[State]map[Event]State)
	for _, transition := range transitions {
		if !slices.Contains(states, transition.State) {
			return nil, ErrUnknownState
		}
		eventStateMap := make(map[Event]State)
		for _, es := range transition.EventStates {
			if !slices.Contains(events, es.Event) {
				return nil, ErrUnknownEvent
			}

			if !slices.Contains(states, es.Next) {
				return nil, ErrUnknownState
			}

			eventStateMap[es.Event] = es.Next
		}
		transitionMap[transition.State] = eventStateMap
	}

	return &FSM{
		current:     start,
		states:      states,
		events:      events,
		transitions: transitionMap,
	}, nil
}

func (fsm *FSM) Current() State {
	fsm.mu.RLock()
	defer fsm.mu.RUnlock()
	return fsm.current
}

func (fsm *FSM) Transition(event Event) (State, error) {
	fsm.mu.Lock()
	// defer fsm.mu.Unlock()
	eventTransition, ok := fsm.transitions[fsm.current]
	if !ok {
		return fsm.current, ErrTerminalState
	}
	nextState, ok := eventTransition[event]
	if !ok {
		return fsm.current, ErrUnknownTransition
	}
	if fsm.current == nextState {
		fsm.mu.Unlock()
		return fsm.current, nil
	}
	fsm.current = nextState
	fsm.mu.Unlock()
	return nextState, nil
}
