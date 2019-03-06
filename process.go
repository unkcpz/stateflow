package giida

import (
  "fmt"
)

// Tasker is a unit run a Execute
type Tasker interface {
  Execute()
}

type Process struct {
  Name string
  task  Tasker
}

func NewProcess(name string, task Tasker) *Process {
  p := &Process{
    Name: name,
    task: task,
  }
  return p
}

func (p *Process) Run() {
  t := p.task
  go func() {
    t.Execute()
    fmt.Print("!!!!")
  }()
}

// type InputGuard struct {
//   ports map[string]bool
//   complete int
// }
//
// //
// func NewInputGuard(ports ...string) *InputGuard {
//   portMap := make(map[string]bool, len(ports))
//   for _, p := range ports {
//     portMap[p] = false
//   }
//   return &InputGuard{portMap, 0}
// }
//
// // Complete is called when a port is closed and returns true when all the ports have been closed
// func (g *InputGuard) Complete(port string) bool {
//   if !g.ports[port] {
//     g.ports[port] = true
//     g.complete++
//   }
//   return g.complete >= len(g.ports)
// }
