package giida

// Process is a unit run a task
type Process interface {
  Task()
}

// Done notifies that the task is finished
type Done struct{}

// Signal is a channel signalling of a completion
type Wait chan struct{}

// Run the Process with Task() function
func Run(c Process) Wait {
  wait := make(Wait)
  go func() {
    c.Task()
    wait <- Done{}
  }()
  return wait
}

type InputGuard struct {
  ports map[string]bool
  complete int
}

//
func NewInputGuard(ports ...string) *InputGuard {
  portMap := make(map[string]bool, len(ports))
  for _, p := range ports {
    portMap[p] = false
  }
  return &InputGuard{portMap, 0}
}

// Complete is called when a port is closed and returns true when all the ports have been closed
func (g *InputGuard) Complete(port string) bool {
  if !g.ports[port] {
    g.ports[port] = true
    g.complete++
  }
  return g.complete >= len(g.ports)
}
