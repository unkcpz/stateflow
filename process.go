package giida

import (
  // "fmt"
  "reflect"
  // "log"
)

type Tasker interface{
  Execute()
}

type Process struct {
  Name string
  task Tasker
  inPorts map[string]chan interface{}
  outPorts map[string]chan interface{}
}

func NewProcess(name string, task Tasker) *Process {
  proc := &Process{
    Name: name,
    task: task,
    inPorts: make(map[string]chan interface{}),
    outPorts: make(map[string]chan interface{}),
  }
  return proc
}

func (p *Process) SetIn(name string, channel chan interface{}) {
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
      v := <- ch
      close(ch)
      val.FieldByName(name).Set(reflect.ValueOf(v))
    }
    task.Execute()
    for name, ch := range p.outPorts {
      ch <- val.FieldByName(name).Interface()
    }
  }()
}
