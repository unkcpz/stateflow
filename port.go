package giida

import (
  // "fmt"
  "reflect"
)

type Tasker interface{
  Execute()
}

type Process struct {
  Name string
  task Tasker
  inPorts map[string]chan int
  outPorts map[string]chan interface{}
}

func NewProcess(name string, task Tasker) *Process {
  proc := &Process{
    Name: name,
    task: task,
    inPorts: make(map[string]chan int),
    outPorts: make(map[string]chan interface{}),
  }
  return proc
}

func (p *Process) SetIn(name string, channel chan int) {
  p.inPorts[name] = channel
}

func (p *Process) SetOut(name string, channel chan interface{}) {
  p.outPorts[name] = channel
}

func (p *Process) Run() {
  task := p.task
  go func() {
    val := reflect.ValueOf(task).Elem()
    for name, ch := range p.inPorts {
      chv := reflect.ValueOf(ch)
      v, _ := chv.Recv()
      val.FieldByName(name).Set(v)
    }
    task.Execute()
    for name, ch := range p.outPorts {
      ch <- val.FieldByName(name).Interface()
    }
  }()
}
