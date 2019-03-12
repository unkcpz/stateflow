package giida

import (
  "reflect"
  "fmt"
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
  val := reflect.ValueOf(task).Elem()
  go func() {
    fmt.Println("00")
    // for name, ch := range p.inPorts {
    //   fmt.Println("01")
    //   fmt.Println(name, p.inPorts)
    //   val.FieldByName(name).Set(reflect.ValueOf(<-ch))
    //   fmt.Println("02")
    //   close(ch)
    //   fmt.Println("03")
    // }
    ch := p.inPorts["X"]
    val.FieldByName("X").Set(reflect.ValueOf(<-ch))
    close(ch)
    ch = p.inPorts["Y"]
    val.FieldByName("Y").Set(reflect.ValueOf(<-ch))
    close(ch)

    task.Execute()
    for name, ch := range p.outPorts {
      ch <- val.FieldByName(name).Interface()
    }
  }()
}
