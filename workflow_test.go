package giida

import (
  "fmt"
  "testing"
)

type PlusOne struct {
  In int
  Out int
}

func (t *PlusOne) Execute() {
  t.Out = t.In + 1
}

func TestSimpleWorkflow(t *testing.T) {
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

    wf := NewWorkflow("test_wf")
    wf.Add(p1)
    wf.Add(p2)
    wf.Connect("p1", "Out", "p2", "In")

    in := make(chan int)
    out := make(chan interface{})
    wf.SetIn("p1", "In", in)
    wf.SetOut("p2", "Out", out)

    fmt.Println("!!!")
    wf.Run()
    in <- test.in

    got := <-out
    if got != test.out {
      t.Errorf("%d + 2 = %d", test.in, got)
    }
  }
}
