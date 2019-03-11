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

type PlusOne struct {
  In int
  Out int
}

func (t *PlusOne) Execute() {
  t.Out = t.In + 1
}

func TestProcess2Process(t *testing.T) {
  tests := []struct {
    in int
    out int
  }{
    {0, 2},
    {-2, 0},
    {199, 201},
  }

  for _, test := range tests {
    p1 := NewProcess("p1", new(PlusOne))
    p2 := NewProcess("p2", new(PlusOne))

    in := make(chan int)
    out := make(chan interface{})
    p1.SetIn("In", in)
    p2.SetOut("Out", out)

    tmpOut := make(chan interface{})
    tmpIn := make(chan int)
    p1.SetOut("Out", tmpOut)
    p2.SetIn("In", tmpIn)
    go func(){
      v := <-tmpOut
      tmpIn <- v.(int)
    }()

    p1.Run()
    p2.Run()

    in <- test.in
    got := <-out
    if got != test.out {
      t.Errorf("%d + 2 = %d", test.in, got)
    }
  }
}
