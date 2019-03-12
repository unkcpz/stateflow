package giida

import (
  "testing"
  "strconv"
)

type AdderToStr struct {
  X int
  Y int
  Sum string
}

func (t *AdderToStr) Execute() {
  t.Sum = strconv.Itoa(t.X + t.Y)
}

func TestProcessWithTwoInputs(t *testing.T) {
  tests := []struct {
    a int
    b int
    sum string
  }{
    {1, 2, "3"},
    {-1, 1, "0"},
  }

  for _, test := range tests {
    proc := NewProcess("adder", new(AdderToStr))

    x := make(chan interface{})
    y := make(chan interface{})
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
