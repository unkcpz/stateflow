package giida

// Process is a unit run a task
type Process interface {
  Task()
}

// Done notifies that the process is finished
type Done struct{}

// Signal is a channel signalling of a completion
type Signal chan struct{}

// Run the Process with Task() function
func Run(c Process) Signal {
  s := make(Signal)
  go func() {
    c.Task()
    s <- Done{}
  }()
  return s
}
