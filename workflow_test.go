package flowmat

import (
  "testing"
)

type PlusOneWF struct {
  In int
  Out int
}

func (t *PlusOneWF) Execute() {
  t.Out = t.In + 1
}

// func TestSimpleWorkflow(t *testing.T) {
//   tests := []struct {
//     in int
//     out int
//   }{
//     {0, 2},
//     {-2, 0},
//     {199, 201},
//   }
//
//   for _, test := range tests {
//     p1 := NewProcess("p1", new(PlusOneWF))
//     p2 := NewProcess("p2", new(PlusOneWF))
//
//     wf := NewWorkflow("test_wf")
//     wf.Add(p1)
//     wf.Add(p2)
//     wf.Connect("p1", "Out", "p2", "In")
//
//     wf.MapIn("wfIn", "p1", "In")
//     wf.MapOut("wfOut", "p2", "Out")
//
//     wf.In("wfIn", test.in)
//     wf.Run()
//
//     got := wf.Out("wfOut")
//     if got.(int) != test.out {
//       t.Errorf("%d + 2 = %d", test.in, got)
//     }
//   }
// }

func TestUseWorkflowAsProcess(t *testing.T) {
  tests := []struct {
    in int
    out int
  }{
    {0, 2},
    {-2, 0},
    {199, 201},
  }

  for _, test := range tests {
    p1 := NewProcess("p1", new(PlusOneWF))
    p2 := NewProcess("p2", new(PlusOneWF))

    wf := NewWorkflow("test_wf")
    wf.Add(p1)
    wf.Add(p2)
    wf.Connect("p1", "Out", "p2", "In")

    wf.MapIn("wfIn", "p1", "In")
    wf.MapOut("wfOut", "p2", "Out")

    in := wf.ExposeIn("wfIn")
    out := wf.ExposeOut("wfOut")

    wf.Run()

    in.channel <-test.in
    got := <-out.channel

    if got.(int) != test.out {
      t.Errorf("%d + 2 = %d", test.in, got)
    }
  }
}

// // // Connect WF->Proc as a new WF
// // func TestWorkflowIsProcess(t *testing.T) {
// //   tests := []struct {
// //     in int
// //     out int
// //   }{
// //     {0, 3},
// //     {-2, 1},
// //     {999, 1002},
// //   }
// //
// //   for _, test := range tests {
// //     p1 := NewProcess("p1", new(PlusOneWF))
// //     p2 := NewProcess("p2", new(PlusOneWF))
// //
// //     subwf := NewWorkflow("wftest")
// //     subwf.Add(p1)
// //     subwf.Add(p2)
// //     subwf.Connect("p1", "Out", "p2", "In")
// //     subwf.ExposeIn("wfIn", "p1", "In")
// //     subwf.ExposeOut("wfOut", "p2", "Out")
// //
// //     pBefore := NewProcess("pb", new(PlusOneWF))
// //     wf := NewWorkflow("wf")
// //     wf.Add(pBefore)
// //     wf.Add(subwf)
// //     wf.Connect("pb", "Out", "wftest", "wfIn")
// //
// //     wf.ExposeIn("wfIn", "pb", "In")
// //     wf.ExposeOut("wfOut", "wftest", "wfOut")
// //
// //     wf.In("wfIn", test.in)
// //     wf.Run()
// //
// //     got := wf.Out("wfOut")
// //     if got.(int) != test.out {
// //       t.Errorf("%d + 3 = %d", test.in, got)
// //     }
// //   }
// // }
