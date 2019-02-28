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
