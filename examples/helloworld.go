package main

import (
	"fmt"
	"github.com/unkcpz/stateflow"
	"log"
	"strings"
)

func main() {
	stateflow.InitLogAudit()

	myname := "Jason"

	wf := stateflow.NewWorkflow("greetWF")
	wf.NewProcess("capin", new(CapIn))
	wf.NewProcess("greet", new(Greet))

	err := wf.Connect("capin", "Out", "greet", "Name")
	if err != nil {
		log.Fatalln(err)
	}
	wf.MapIn("wfIn", "capin", "In")
	wf.MapOut("wfOut", "greet", "Greeting")

	wf.In("wfIn", myname)
	wf.Flow()

	greeting := wf.Out("wfOut")
	fmt.Println(greeting)
}

// Task to capitalize In
type CapIn struct {
	In  string
	Out string
}

func (t *CapIn) Execute() {
	t.Out = strings.ToUpper(t.In)
}

// Task for greeting
type Greet struct {
	Name     string
	Greeting string
}

func (t *Greet) Execute() {
	t.Greeting = fmt.Sprintf("Hello %s.", t.Name)
}
