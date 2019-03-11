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
    // proc.BindPort("pIn", "In")
    // proc.BindPort("pOut", "Out")

    in := make(chan int)
    out := make(chan string)
    proc.SetIn("pIn", in)
    proc.SetOut("pOut", out)

    proc.Run()
    in <- test.in

    got := <-out
    if got != test.expected {
      t.Errorf("string(%d+1) = %s, expected %s", test.in, got, test.expected)
    }
  }
}
