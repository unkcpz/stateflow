package main

import (
  "strings"
  "fmt"
	"log"
  "github.com/unkcpz/gflow"
)

func main() {
	gflow.InitLogAudit()

  myname := "Jason"
  p1 := gflow.NewProcess("capin", new(CapIn))
  p2 := gflow.NewProcess("greet", new(Greet))

  wf := gflow.NewWorkflow("greetWF")
  wf.Add(p1)
  wf.Add(p2)

  err := wf.Connect("capin", "Out", "Greet", "Name")
  if err != nil {
    log.Println(err)
  }
  wf.MapIn("wfIn", "capin", "In")
  wf.MapOut("wfOut", "greet", "Greeting")

  wf.In("wfIn", myname)
  wf.Load()
  wf.Start()
  wf.Finish()

  greeting := wf.Out("wfOut")
  fmt.Println(greeting)
}

// Task to capitalize In
type CapIn struct {
  In string
  Out string
}

func (t *CapIn) Execute() {
  t.Out = strings.ToUpper(t.In)
}

// Task for greeting
type Greet struct {
  Name string
  Greeting string
}

func (t *Greet) Execute() {
  t.Greeting = fmt.Sprintf("Hello %s.", t.Name)
}
