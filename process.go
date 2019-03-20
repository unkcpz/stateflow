package flowmat

import (
  "reflect"
  "sync"
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
  go func(){
    task := p.task
    val := reflect.ValueOf(task).Elem()
    var wg sync.WaitGroup
    for name, ch := range p.inPorts {
      wg.Add(1)
      go func(name string, ch chan interface{}) {
        defer wg.Done()
        val.FieldByName(name).Set(reflect.ValueOf(<-ch))
        close(ch)
      }(name, ch)
    }
    wg.Wait()
    task.Execute()
    for name, ch := range p.outPorts {
      wg.Add(1)
      go func(name string, ch chan interface{}) {
        defer wg.Done()
        ch <- val.FieldByName(name).Interface()
        close(ch)
      }(name, ch)
    }
    wg.Wait()
  }()
}
