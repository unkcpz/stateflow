package flowmat

import (
  "testing"
  "strconv"
  "log"
)

func TestComplexParaWorkflow(t *testing.T) {
  tests := []struct {
    in string
    out float64
  }{
    // {"0", 0.0},
    {"0.1", 0.8},
    // {"199", 1592},
  }

  for _, test := range tests {
    conv := NewProcess("conv", new(StringToFloat))
    t1 := NewProcess("time1", new(TimeTwo))
    t2 := NewProcess("time2", new(TimeTwo))
    t3 := NewProcess("time3", new(TimeTwo))
    t4 := NewProcess("time4", new(TimeTwo))
    sum := NewProcess("sumAll", new(SumAll))
    clone := NewProcess("clone4", new(CloneFour))

    wf := NewWorkflow("ComplexParaWF")
    wf.Add(conv)
    wf.Add(t1)
    wf.Add(t2)
    wf.Add(t3)
    wf.Add(t4)
    wf.Add(sum)
    wf.Add(clone)

    wf.Connect("conv", "Out", "clone4", "In")
    wf.Connect("clone4", "Out0", "time1", "In")
    wf.Connect("clone4", "Out1", "time2", "In")
    wf.Connect("clone4", "Out2", "time3", "In")
    wf.Connect("clone4", "Out3", "time4", "In")
    wf.Connect("time1", "Out", "sumAll", "In0")
    wf.Connect("time2", "Out", "sumAll", "In1")
    wf.Connect("time3", "Out", "sumAll", "In2")
    wf.Connect("time4", "Out", "sumAll", "In3")

    wf.MapIn("wfIn", "conv", "In")
    wf.MapOut("wfOut", "sumAll", "Out")

    wf.In("wfIn", test.in)
    wf.Load()
    wf.Start()
    wf.Finish()

    got := wf.Out("wfOut")
    if got.(float64) != test.out {
      t.Errorf("float(%s) * 8 = %f", test.in, got)
    }
  }
}

// //
type TimeTwo struct {
  In float64
  Out float64
}

func (t *TimeTwo) Execute() {
  t.Out = t.In * 2
}

// //
type StringToFloat struct {
  In string
  Out float64
}

func (t *StringToFloat) Execute() {
  var err error
  t.Out, err = strconv.ParseFloat(t.In, 64)
  if err != nil {
    log.Fatalf("Executing with task: %v, failed", t)
  }
}

// //
type CloneFour struct {
  In float64
  Out0 float64
  Out1 float64
  Out2 float64
  Out3 float64
}

func (t *CloneFour) Execute() {
  t.Out0 = t.In
  t.Out1 = t.In
  t.Out2 = t.In
  t.Out3 = t.In
}

// //
type SumAll struct {
  In0 float64
  In1 float64
  In2 float64
  In3 float64
  Out float64
}

func (t *SumAll) Execute() {
  t.Out = t.In0 + t.In1 + t.In2 + t.In3
}
