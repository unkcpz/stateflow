package giida

import (
  "testing"
)

func TestSetWfName(t *testing.T) {
  wf := NewWorkflow("workflowname", 4)

  expectedWfName := "workflowname"
  if wf.name != expectedWfName {
    t.Errorf("Workflow name is wrong, got %s, expect %s\n", wf.name, expectedWfName)
  }
}

func TestMaxGoroutineCapacity(t *testing.T) {
  wf := NewWorkflow("tmp", 16)

  if got := cap(wf.goroutineTasks); got != 16 {
    t.Errorf("Workflow goroutine cap number is %d, expect 16", got)
  }
}
