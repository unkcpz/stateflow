package giida

import (
  "testing"
)

func TestSetWfName(t *testing.T) {
  wf := NewWorkflow("workflowname", 8, 0)

  expectedWfName := "workflowname"
  if wf.name != expectedWfName {
    t.Errorf("Workflow name is wrong, got %s, expect %s\n", wf.name, expectedWfName)
  }
}

func newDoubleAddOneEcho() (*Workflow, error) {
  n := NewWorkflow("new", 8, 0)
  // Task
  e1 := NewProcess("e1", new(echo))
  e2 := NewProcess("e2", new(echo))

  if err := n.Add("e1", *e1); err != nil {
    return nil, err
  }
  if err := n.Add("e2", *e2); err != nil {
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

func testWorkflowWithNumberSequence(t *testing.T) {
  tests := []struct {
    in int
    expected int
  } {
    {93, 95},
    {52, 54},
    {1, 3},
    {24, 26},
    {35, 37},
  }

  for _, test := range tests {
    n, err := newDoubleAddOneEcho()
    if err != nil {
      t.Error(err)
      return
    }
    
    in := make(chan int)
    out := make(chan int)
    n.SetInPort("netIn", in)
    n.SetOutPort("netOut", out)

    wait := n.Run()

    in <- test.in
    if got := <-out; got != test.expected {
      t.Errorf("%d + 2 != %d", test.in, test.expected)
    }
    <-wait
  }
}

// Task for test
type echo struct {
  In <-chan int
  Out chan<- int
}

func (c *echo) Execute() {
  in := <-c.In
  c.Out <- in + 1
}
