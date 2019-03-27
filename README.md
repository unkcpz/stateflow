# stateflow
Pure Go Automated Interactive Infrastructure and Database for Computational Science

## Introduction

## Examples

### Hello World

```go
package main

import (
  "strings"
  "fmt"
	"log"
  "github.com/unkcpz/stateflow"
)

func main() {
	stateflow.InitLogAudit()

  myname := "Jason"
  p1 := stateflow.NewProcess("capin", new(CapIn))
  p2 := stateflow.NewProcess("greet", new(Greet))

  wf := stateflow.NewWorkflow("greetWF")
  wf.Add(p1)
  wf.Add(p2)

  err := wf.Connect("capin", "Out", "greet", "Name")
  if err != nil {
    log.Fatalln(err)
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
```

### Simple Process

## Simple Workflow

## Acknowledgements

<!-- - stateflow is very heavily dependent on the proven principles form [Flow-Based
  Programming (FBP)](http://www.jpaulmorrison.com/fbp), as invented by [John Paul Morrison](http://www.jpaulmorrison.com/fbp).
  From Flow-based programming, stateflow uses the ideas of separate network
  (workflow dependency graph) definition, named in- and out-ports,
  sub-networks/sub-workflows and bounded buffers (already available in Go's
  channels) to make writing workflows as easy as possible. -->
- This library is has been much influenced/inspired also by the
  [GoFlow](https://github.com/trustmaster/goflow) library by [Vladimir Sibirov](https://github.com/trustmaster/goflow),
  and [SciPipe](https://github.com/scipipe/scipipe) library by [Samuel Lampa](https://github.com/samuell)
  and [AiiDA](http://www.aiida.net/) library by [AiiDA team](http://www.aiida.net/team/)
