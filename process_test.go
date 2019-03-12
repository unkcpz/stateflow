package flowmat

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

type MyType struct {
  Adder []int
  Scaler int
}

type ComplexTask struct {
  MyT MyType
  Inc int
  Out string
}

func (t *ComplexTask) Execute() {
  myType := t.MyT
  sum := 0
  for _, i := range myType.Adder {
    sum += i
  }
  sum *= myType.Scaler
  t.Out = strconv.Itoa(sum + t.Inc)
}

func TestComplexProcessWithCustomType(t *testing.T) {
  tests := []struct {
    mt MyType
    inc int
    expected string
  }{
    {
      MyType{[]int{1, 2}, 2},
      100,
      "106",
    },
    {
      MyType{[]int{1, 2, 4, 6}, 0},
      1000,
      "1000",
    },
  }

  for _, test := range tests {
    proc := NewProcess("complexTask", new(ComplexTask))

    myType := make(chan interface{})
    inc := make(chan interface{})
    out := make(chan interface{})
    proc.SetIn("MyT", myType)
    proc.SetIn("Inc", inc)
    proc.SetOut("Out", out)

    proc.Run()
    myType <- test.mt
    inc <- test.inc

    got := <-out
    if got != test.expected {
      t.Errorf("process ComplexTask[MyT=%v, Inc=%d], got %s, expected %s", test.mt, test.inc, got, test.expected)
    }
  }
}
