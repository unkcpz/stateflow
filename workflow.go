package giida

type Workflow struct {
  name string
  goroutineTasks  chan struct{}
}

func NewWorkflow(name string, maxGoroutineTasks int) *Workflow {
  wf := &Workflow{
    name: name,
    goroutineTasks: make(chan struct{}, maxGoroutineTasks),
  }
  return wf
}
