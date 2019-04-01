# stateflow
Pure Go Automated Interactive Infrastructure and Database for Computational Science

## Install

```sh
▶ go get -u github.com/unkcpz/stateflow
```

## Introduction

## Examples

### Hello World (Simple Workflow)

![](https://raw.githubusercontent.com/unkcpz/images/master/stateflow.repo/README-simpleWF.png)
```go
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/unkcpz/stateflow"
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
```

### Simple Process
![](https://raw.githubusercontent.com/unkcpz/images/master/stateflow.repo/README-2in2out.png)
```go
package main

import (
	"fmt"

	"github.com/unkcpz/stateflow"
)

func main() {
	num := 10
	deno := 3

	stateflow.InitLogAudit()
	proc := stateflow.NewProcess("2in2out", new(TwoResults))
	proc.In("Num", num)
	proc.In("Deno", deno)

	proc.Flow()

	quot := proc.Out("Quot")
	// Output can be not export and used
	rem := proc.Out("Rem")

	fmt.Printf("%d / %d = (Quot: %d, Rem %d)\n", num, deno, quot, rem)
}

type TwoResults struct {
	Num  int
	Deno int
	Quot int
	Rem  int
}

func (t *TwoResults) Execute() {
	t.Quot = t.Num / t.Deno
	t.Rem = t.Num % t.Deno
}
```

## SubWorkflow
![](https://raw.githubusercontent.com/unkcpz/images/master/stateflow.repo/README-subwf.png)
```go
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
```

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
