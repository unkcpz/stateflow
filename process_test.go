package flowmat

import (
  "testing"
  "strconv"
  // "fmt"
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

    x := proc.ExposeIn("X")
    y := proc.ExposeIn("Y")
    sum := proc.ExposeOut("Sum")

    proc.Load()
    x.Feed(test.a)
    y.Feed(test.b)

    got := sum.Extract()

    if got !=test.sum {
      t.Errorf("component: %d + %d == %d", test.a, test.b, got)
    }
  }
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

    num := proc.ExposeIn("Num")
    deno := proc.ExposeIn("Deno")
    quot := proc.ExposeOut("Quot")
    rem := proc.ExposeOut("Rem")

    proc.Load()
    num.Feed(test.num)
    deno.Feed(test.deno)

    gQuot := quot.Extract()
    gRem := rem.Extract()
    if gQuot != test.expectQuot || gRem != test.expectRem {
      t.Errorf("%d / %d = (Quot: %d, Rem: %d)", test.num, test.deno, gQuot, gRem)
    }
  }
}

// Test process of multiple output int, int -> int, int
func TestProcessTwoInTwoOutNotAllUsed(t *testing.T) {
  tests := []struct {
    num int
    deno int
    expectQuot int
    expectRem int
  }{
    {5, 4, 1, 1},
    {10, 3, 3, 1},
  }

  // What if one output is not exposed and used?
  for _, test := range tests {
    proc := NewProcess("2to2", new(TwoResults))

    num := proc.ExposeIn("Num")
    deno := proc.ExposeIn("Deno")
    quot := proc.ExposeOut("Quot")

    proc.Load()
    num.Feed(test.num)
    deno.Feed(test.deno)
    proc.Finish()

    gQuot := quot.Extract()
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

    myType := proc.ExposeIn("MyT")
    inc := proc.ExposeIn("Inc")
    out := proc.ExposeOut("Out")

    proc.Load()
    myType.Feed(test.mt)
    inc.Feed(test.inc)

    got := out.Extract()
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

// Test Process as independent plugin two int input and string output
func TestProcessWithTwoInputsPlugin(t *testing.T) {
  tests := []struct {
    a int
    b int
    sum string
  }{
    {1, 2, "3"},
    {-1, 1, "0"},
  }

  for _, test := range tests {
    proc := NewProcess("adder", new(AdderToStrPlugin))

    proc.In("X", test.a)
    proc.In("Y", test.b)

    proc.Load()
    proc.Start()
    proc.Finish()

    got := proc.Out("Sum")
    if got !=test.sum {
      t.Errorf("plugin: %d + %d == %s", test.a, test.b, got)
    }
  }
}


type AdderToStrPlugin struct {
  X int
  Y int
  Sum string
}

func (t *AdderToStrPlugin) Execute() {
  t.Sum = strconv.Itoa(t.X + t.Y)
}
