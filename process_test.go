package flowmat

import (
  "testing"
  "strconv"
)

// Test Process two int input and string output
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
      t.Errorf("component: %d + %d == %d", test.a, test.b, got)
    }
  }

  // // Test Process used as an independent component
  // for _, test := range tests {
  //   proc := NewProcess("adder", new(AdderToStr))
  //
  //   proc.In("X", test.a)
  //   proc.In("Y", test.b)
  //
  //   proc.Run()
  //
  //   got := proc.Out("Sum")
  //   if got !=test.sum {
  //     t.Errorf("plugin: %d + %d == %d", test.a, test.b, got)
  //   }
  // }
}


type AdderToStr struct {
  X int
  Y int
  Sum string
}

func (t *AdderToStr) Execute() {
  t.Sum = strconv.Itoa(t.X + t.Y)
}

// Test process of multiple output int, int -> int, int
func TestProcessTwoInTwoOut(t *testing.T) {
  tests := []struct {
    num int
    deno int
    expectQuot int
    expectRem int
  }{
    {5, 4, 1, 1},
    {10, 3, 3, 1},
  }

  for _, test := range tests {
    proc := NewProcess("2to2", new(TwoResults))

    num := make(chan interface{})
    deno := make(chan interface{})
    quot := make(chan interface{})
    rem := make(chan interface{})
    proc.SetIn("Num", num)
    proc.SetIn("Deno", deno)
    proc.SetOut("Quot", quot)
    proc.SetOut("Rem", rem)

    proc.Run()
    num <- test.num
    deno <- test.deno

    gQuot := <-quot
    gRem := <-rem
    if gQuot != test.expectQuot || gRem != test.expectRem {
      t.Errorf("%d / %d = (Quot: %d, Rem: %d)", test.num, test.deno, gQuot, gRem)
    }
  }

  // What if one output is not used?
  for _, test := range tests {
    proc := NewProcess("2to2", new(TwoResults))

    num := make(chan interface{})
    deno := make(chan interface{})
    quot := make(chan interface{})
    rem := make(chan interface{})
    proc.SetIn("Num", num)
    proc.SetIn("Deno", deno)
    proc.SetOut("Quot", quot)
    proc.SetOut("Rem", rem)

    proc.Run()
    num <- test.num
    deno <- test.deno

    gQuot := <-quot
    // gRem := <-rem
    <-rem
    if gQuot != test.expectQuot{
      t.Errorf("%d / %d = (Quot: %d)", test.num, test.deno, gQuot)
    }
  }
}

type TwoResults struct {
  Num int
  Deno int
  Quot int
  Rem int
}

func (t *TwoResults) Execute() {
  t.Quot = t.Num / t.Deno
  t.Rem = t.Num % t.Deno
}

// Test a complex task with multi operation and custom defined type
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
