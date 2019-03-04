package giida

import (
  "testing"
)

func TestSetWfName(t *testing.T) {
  wf := NewWorkflow("workflowname", 8)

  expectedWfName := "workflowname"
  if wf.name != expectedWfName {
    t.Errorf("Workflow name is wrong, got %s, expect %s\n", wf.name, expectedWfName)
  }
}

func newDoubleEcho() (*Workflow, error) {
  n := NewWorkflow()
  // Process
  e1 := new(echo)
  e2 := new(echo)

  if err := n.Add("e1", e1); err != nil {
    return nil, err
  }
  if err := n.Add("e2", e2); err != nil {
    return nil, err
  }
  if err := n.Connect("e1", "Out", "e2", "In"); err != nil {
    return nil, err
  }
  // Ports
  if err := n.MapInPort("netIn", "e1", "In"); err != nil {
    return nil, err
  }
  if err := n.MapOutPort("netOut", "e2", "Out"); err != nil {
    return nil, err
  }

  return n, nil
}

func TestSimpleWorkflow(t *testing.T) {
  n, err := newDoubleEcho()
  if err != nil {
    t.Error(err)
    return
  }

  testWorkflowWithNumberSequence(n, t)
}

func testWorkflowWithNumberSequence(n *Workflow, t *testing.T) {
  data := []int{93, 52, 1, 24, 35, 63, 634, 12}

  in := make(chan int)
  out := make(chan int)
  n.SetInPort("netIn", in)
  n.SetOutPort("netOut", out)

  wait := Run(n)

  go func() {
    for _, n := range data {
      in <- n
    }
    close(in)
  }()

  i := 0
  for got := range out {
    expected := data[i]
    if got != expected {
      t.Errorf("%d != %d\n", got, expected)
    }
    i++
  }

  <- wait
}

// Process for test
type echo struct {
  In <-chan int
  Out chan<- int
}

func (c *echo) Task() {
  for i := range c.In {
    c.Out <- i
  }
}
