package giida

// Tasker is a unit run a Execute
type Tasker interface {
  Execute()
}

// Done notifies that the Execute is finished
type Done struct{}

// Signal is a channel signalling of a completion
type Wait chan struct{}

// Run the Tasker with Execute() function
func Run(c Tasker) Wait {
  wait := make(Wait)
  go func() {
    c.Execute()
    wait <- Done{}
  }()
  return wait
}

type Process struct {
  Name string
  Task  Tasker
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
