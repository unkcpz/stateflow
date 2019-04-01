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
