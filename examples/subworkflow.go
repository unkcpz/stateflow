package main

import (
	"fmt"
	"strconv"

	"github.com/unkcpz/stateflow"
)

func main() {
	stateflow.InitLogAudit()
	// define sub-workflow "doubleplus"
	subwf := stateflow.NewWorkflow("doubleplus")
	subwf.NewProcess("p1", new(PlusOne))
	subwf.NewProcess("p2", new(PlusOne))
	subwf.Connect("p1", "Out", "p2", "In")
	subwf.MapIn("wfIn", "p1", "In")
	subwf.MapOut("wfOut", "p2", "Out")

	// define process link after the sub-workflow
	pAfter := stateflow.NewProcess("pa", new(IntToString))

	// define main Workflow
	wf := stateflow.NewWorkflow("mainWF")
	wf.Add(pAfter)
	wf.Add(subwf)
	wf.Connect("doubleplus", "wfOut", "pa", "In")
	wf.MapIn("wfIn", "doubleplus", "wfIn")
	wf.MapOut("wfOut", "pa", "Out")

	in := 1
	wf.In("wfIn", in)
	wf.Flow()

	fmt.Printf("%d + 2 = %s\n", in, wf.Out("wfOut"))
}

type PlusOne struct {
	In  int
	Out int
}

func (t *PlusOne) Execute() {
	t.Out = t.In + 1
}

type IntToString struct {
	In  int
	Out string
}

func (t *IntToString) Execute() {
	t.Out = strconv.Itoa(t.In)
}
