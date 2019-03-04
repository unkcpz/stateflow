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

func TestProcessWithTwoInputs(t *testing.T) {
  tests := []struct {
    op1 int
    op2 int
    sums int
  }{
    {3, 38, 41},
    {3, 4, 7},
    {92, 4, 96},
  }

  in1 := make(chan int)
  in2 := make(chan int)
  out := make(chan int)
  c := &adder{in1, in2, out}

  wait := Run(c)

  go func() {
    for _, t := range tests {
      in1 <-t.op1
    }
    close(in1)
  }()
  go func() {
    for _, t := range tests {
      in2 <-t.op2
    }
    close(in2)
  }()

  for _, test := range tests {
    got := <-out
    expected := test.sums
    if got != expected {
      t.Errorf("%d + %d = %d, expect %d", test.op1, test.op2, got, expected)
    }
  }

  <-wait
}

type adder struct {
  Op1 <-chan int
  Op2 <-chan int
  Sum chan<- int
}

func (c *adder) Task() {
  guard := NewInputGuard("op1", "op2")

  op1Buf := make([]int, 0, 10)
  op2Buf := make([]int, 0, 10)
  addOp := func(op int, buf, otherBuf *[]int) {
    if len(*otherBuf) > 0 {
      otherOp := (*otherBuf)[0]
      *otherBuf = (*otherBuf)[1:]
      c.Sum <- (op + otherOp)
    } else {
      *buf = append(*buf, op)
    }
  }

  for {
    select {
    case op1, ok := <-c.Op1:
      if ok {
        addOp(op1, &op1Buf, &op2Buf)
      } else if guard.Complete("op1") {
        return
      }
    case op2, ok := <-c.Op2:
      if ok {
        addOp(op2, &op2Buf, &op1Buf)
      } else if guard.Complete("op2") {
        return
      }
    }
  }
}
