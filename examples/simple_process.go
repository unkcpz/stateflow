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
