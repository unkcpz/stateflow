package giida

import (
  "testing"
  "fmt"
)

type Demo struct {
  In int
  Out string
}

func (t *Demo) Execute() {
  i := t.In
  t.Out = fmt.Sprintf("%d", i+1)
}

func TestProcess(t *testing.T) {
  tests := []struct{
    in int
    expected string
  }{
    {5, "6"},
    {10, "11"},
    {-1, "0"},
    {3, "4"},
  }

  for _, test := range tests {
    proc := NewProcess("test", new(Demo))

    in := make(chan int)
    out := make(chan interface{})
    proc.SetIn("In", in)
    proc.SetOut("Out", out)

    proc.Run()
    in <- test.in

    got := <-out
    if got != test.expected {
      t.Errorf("string(%d+1) = %s, expected %s", test.in, got, test.expected)
    }
  }
}

type Adder struct {
  X int
  Y int
  Sum int
}

func (t *Adder) Execute() {
  t.Sum = t.X + t.Y
}

func TestProcessWithTwoInputs(t *testing.T) {
  tests := []struct {
    a int
    b int
    sum int
  }{
    {1, 2, 3},
    {-1, 1, 0},
  }

  for _, test := range tests {
    proc := NewProcess("adder", new(Adder))

    x := make(chan int)
    y := make(chan int)
    sum := make(chan interface{})
    proc.SetIn("X", x)
    proc.SetIn("Y", y)
    proc.SetOut("Sum", sum)

    proc.Run()
    x <- test.a
    y <- test.b

    got := <-sum
    if got !=test.sum {
      t.Errorf("%d + %d == %d", test.a, test.b, got)
    }
  }
}
