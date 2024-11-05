package fsm_test

import (
	"testing"

	"github.com/bento01dev/fsm"
)

const (
	STATE_1 fsm.State = iota + 1
	STATE_2
	STATE_3
	STATE_4
	STATE_5
)

const (
	EVENT_1 fsm.Event = iota + 1
	EVENT_2
	EVENT_3
	EVENT_4
	EVENT_5
	EVENT_6
)

func BenchmarkFSM(b *testing.B) {
	states := []fsm.State{STATE_1, STATE_2, STATE_3, STATE_4, STATE_5}
	events := []fsm.Event{EVENT_1, EVENT_2, EVENT_3, EVENT_4, EVENT_5, EVENT_6}
	transitions := []fsm.Transition{
		{
			State: STATE_1,
			EventStates: []fsm.EventState{
				{
					Event: EVENT_2,
					Next:  STATE_2,
				},
			},
		},
		{
			State: STATE_2,
			EventStates: []fsm.EventState{
				{
					Event: EVENT_1,
					Next:  STATE_1,
				},
			},
		},
	}
	f, err := fsm.NewFSM(states, events, STATE_1, transitions...)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := f.Transition(EVENT_2); err != nil {
			b.Fatal(err)
		}
		if _, err := f.Transition(EVENT_1); err != nil {
			b.Fatal(err)
		}
	}
}
