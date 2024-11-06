package fsm_test

import (
	"context"
	"testing"

	"github.com/bento01dev/fsm"
	csfsm "github.com/cocoonspace/fsm"
	llfsm "github.com/looplab/fsm"
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
	f := llfsm.NewFSM(
		"foo",
		llfsm.Events{
			{Name: "foo", Src: []string{"foo"}, Dst: "bar"},
			{Name: "bar", Src: []string{"bar"}, Dst: "foo"},
		},
		llfsm.Callbacks{},
	)
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			f.Event(context.Background(), "foo")
			f.Event(context.Background(), "bar")
		}
	})
}

func BenchmarkMyFSM(b *testing.B) {
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
				{
					Event: EVENT_1,
					Next:  STATE_1,
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
				{
					Event: EVENT_2,
					Next:  STATE_2,
				},
			},
		},
	}
	f, err := fsm.NewFSM(states, events, STATE_1, transitions...)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			if _, err := f.Transition(EVENT_2); err != nil {
				b.Fatal(err)
			}
			if _, err := f.Transition(EVENT_1); err != nil {
				b.Fatal(err)
			}
		}
	})
}

const (
	StateFoo csfsm.State = iota
	StateBar
)

const (
	EventFoo csfsm.Event = iota
	EventBar
)

func BenchmarkCocoonSpaceFSM(b *testing.B) {
	f := csfsm.New(StateFoo)
	f.Transition(csfsm.On(EventFoo), csfsm.Src(StateFoo), csfsm.Dst(StateBar))
	f.Transition(csfsm.On(EventBar), csfsm.Src(StateBar), csfsm.Dst(StateFoo))
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			f.Event(EventFoo)
			f.Event(EventBar)
		}
	})
}

func BenchmarkLooplabFSM(b *testing.B) {
	f := llfsm.NewFSM(
		"foo",
		llfsm.Events{
			{Name: "foo", Src: []string{"foo"}, Dst: "bar"},
			{Name: "bar", Src: []string{"bar"}, Dst: "foo"},
		},
		llfsm.Callbacks{},
	)
	b.ResetTimer()
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			f.Event(context.Background(), "foo")
			f.Event(context.Background(), "bar")
		}
	})
}
