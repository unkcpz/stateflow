package giida

import (
  // "fmt"
)

// type Tasker interface{
//   Execute()
// }

type Process struct {
  Name string
  task *Demo
  inPorts map[string]chan int
  outPorts map[string]chan string
}

func NewProcess(name string, task *Demo) *Process {
  proc := &Process{
    Name: name,
    task: task,
    inPorts: make(map[string]chan int),
    outPorts: make(map[string]chan string),
  }
  return proc
}

// func (p *Process) BindPort(customName, taskName string) {
//
// }

func (p *Process) SetIn(name string, channel chan int) {
  p.inPorts[name] = channel
}

func (p *Process) SetOut(name string, channel chan string) {
  p.outPorts[name] = channel
}

func (p *Process) Run() {
  task := p.task
  go func() {
    task.In = <-p.inPorts["pIn"]
    task.Execute()
    p.outPorts["pOut"] <- task.Out
  }()
}
