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
