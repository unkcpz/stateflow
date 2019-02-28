package giida

import (
  "testing"
)

// Test a simple Process that runs only once
func TestSimpleProcess(t *testing.T) {
  in := make(chan int)
  out := make(chan int)
  c := &doubleOnce{
    in,
    out,
  }

  wait := Run(c)

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

func (c *doubleOnce) Task() {
  i := <-c.In
  c.Out <- 2 * i
}
