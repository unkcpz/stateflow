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
