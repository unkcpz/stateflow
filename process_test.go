package giida

import (
  "testing"
)

// Test a simple Process that runs only once
func TestSimpleProcess(t *testing.T) {
  in := make(chan int)
  out := make(chan int)
  proc := &doubleOnce{
    in,
    out,
  }

  wait := Run(proc)

  in <- 21
  if got := <-out; got != 42 {
    t.Errorf("%d != %d", got, 42)
  }

  <-wait
}

type doubleOnce struct {
  In <-chan int
  Out chan<- int
}

func (proc *doubleOnce) Task() {
  i := <-proc.In
  proc.Out <- 2 * i
}

// Test a simple long running process with one input
func TestSimpleLongRunningProcess(t *testing.T) {
  tests := []struct {
    in int
    expected int
  }{
    {12, 24},
    {5, 10},
  }

  in := make(chan int)
  out := make(chan int)
  proc := &doubler{
    in,
    out,
  }

  wait := Run(proc)

  for _, test := range tests {
    in <- test.in
    got := <-out

    if got != test.expected {
      t.Errorf("%d != %d", got, test.expected)
    }
  }

  close(in)
  <-wait
}

type doubler struct {
  In <-chan int
  Out chan<- int
}

func (proc *doubler) Task() {
  for i := range proc.In {
    proc.Out <- 2 * i
  }
}
