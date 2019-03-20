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
  inPorts map[string]*Port
  outPorts map[string]*Port
}

func NewProcess(name string, task Tasker) *Process {
  proc := &Process{
    Name: name,
    task: task,
    inPorts: make(map[string]*Port),
    outPorts: make(map[string]*Port),
  }
  return proc
}

func (p *Process) SetIn(name string, channel chan interface{}) {
  p.inPorts[name] = &Port{
    channel: channel,
  }
}

func (p *Process) SetOut(name string, channel chan interface{}) {
  p.outPorts[name] = &Port{
    channel: channel,
  }
}

func (p *Process) Run() {
  go func(){
    task := p.task
    val := reflect.ValueOf(task).Elem()
    var wg sync.WaitGroup
    for name, port := range p.inPorts {
      wg.Add(1)
      go func(name string, port *Port) {
        defer wg.Done()
        ch := port.channel
        v := reflect.ValueOf(<-ch)
        val.FieldByName(name).Set(v)
        if port.cache == nil{
          port.cache = v.Interface()
        }
        close(ch)
      }(name, port)
    }
    wg.Wait()

    // Execute the function of Process
    task.Execute()

    for name, port := range p.outPorts {
      wg.Add(1)
      go func(name string, port *Port) {
        defer wg.Done()
        ch := port.channel
        v := val.FieldByName(name).Interface()
        ch <- v
        port.cache = v
        close(ch)
      }(name, port)
    }
    wg.Wait()
  }()
}
