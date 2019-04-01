package fsm

import (
  "testing"
)

func TestSameState(t *testing.T) {
  fsm := NewFSM(
    "start",
    Events{
      {Name: "run", Src: []string{"start"}, Dst: "start"},
    },
    Callbacks{},
  )
  fsm.Event("run")
  if fsm.Current() != "start" {
    t.Error("expected state to be 'start'")
  }
}

func TestSetState(t *testing.T) {
  fsm := NewFSM(
    "walking",
    Events{
      {Name: "walk", Src: []string{"start"}, Dst: "walking"},
    },
    Callbacks{},
  )
  fsm.SetState("start")
  if fsm.Current() != "start" {
    t.Error("expected state to be 'start'")
  }
  err := fsm.Event("walk")
  if err != nil {
    t.Error("transition is expected no error")
  }
}

type fakeTransitionerObj struct {
}

func (t fakeTransitionerObj) transition(f *FSM) error {
  return &InternalError{}
}

func TestBadTransition(t *testing.T) {
  fsm := NewFSM(
    "start",
    Events{
      {Name: "run", Src: []string{"start"}, Dst: "running"},
    },
    Callbacks{},
  )
  fsm.transitionerObj = new(fakeTransitionerObj)
  err := fsm.Event("run")
  if err == nil {
    t.Error("bad transition should give an error")
  }
}

func TestInappropriateEvent(t *testing.T) {
  fsm := NewFSM(
    "closed",
    Events{
      {Name: "open", Src: []string{"closed"}, Dst: "open"},
      {Name: "close", Src: []string{"open"}, Dst: "closed"},
    },
    Callbacks{},
  )
  err := fsm.Event("close")
  if e, ok := err.(InvalidEventError); !ok && e.Event != "close" && e.State != "closed" {
    t.Error("expected 'InvalidEventError' with correct state and event")
  }
}

func TestInvaliedEvent(t *testing.T) {
  fsm := NewFSM(
    "closed",
    Events{
      {Name: "open", Src: []string{"closed"}, Dst: "open"},
      {Name: "close", Src: []string{"open"}, Dst: "closed"},
    },
    Callbacks{},
  )
  err := fsm.Event("lock")
  if e, ok := err.(UnknownEventError); !ok && e.Event != "close" {
    t.Error("expected 'UnknownEventError' with correct event")
  }
}

func TestMultipleSources(t *testing.T) {
  fsm := NewFSM(
    "one",
    Events{
      {Name: "first", Src: []string{"one"}, Dst: "two"},
      {Name: "second", Src: []string{"two"}, Dst: "three"},
      {Name: "reset", Src: []string{"one", "two", "three"}, Dst: "one"},
    },
    Callbacks{},
  )

  fsm.Event("first")
  if fsm.Current() != "two" {
    t.Error("expected state to be 'two'")
  }
  fsm.Event("reset")
  if fsm.Current() != "one" {
    t.Error("expected state to be 'one'")
  }
  fsm.Event("first")
  fsm.Event("second")
  if fsm.Current() != "three" {
    t.Error("expected state to be 'three'")
  }
  fsm.Event("reset")
  if fsm.Current() != "one" {
    t.Error("expected state to be 'one'")
  }
}

func TestMultipleEvents(t *testing.T) {
	fsm := NewFSM(
		"start",
		Events{
			{Name: "first", Src: []string{"start"}, Dst: "one"},
			{Name: "second", Src: []string{"start"}, Dst: "two"},
			{Name: "reset", Src: []string{"one"}, Dst: "reset_one"},
			{Name: "reset", Src: []string{"two"}, Dst: "reset_two"},
			{Name: "reset", Src: []string{"reset_one", "reset_two"}, Dst: "start"},
		},
		Callbacks{},
	)

	fsm.Event("first")
	fsm.Event("reset")
	if fsm.Current() != "reset_one" {
		t.Error("expected state to be 'reset_one'")
	}
	fsm.Event("reset")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}

	fsm.Event("second")
	fsm.Event("reset")
	if fsm.Current() != "reset_two" {
		t.Error("expected state to be 'reset_two'")
	}
	fsm.Event("reset")
	if fsm.Current() != "start" {
		t.Error("expected state to be 'start'")
	}
}
