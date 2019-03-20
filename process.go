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
  inPorts map[string]*port
  outPorts map[string]*port
}

func NewProcess(name string, task Tasker) *Process {
  proc := &Process{
    Name: name,
    task: task,
    inPorts: make(map[string]*port),
    outPorts: make(map[string]*port),
  }
  return proc
}

func (p *Process) SetIn(name string, channel chan interface{}) {
  p.inPorts[name] = &port{
    channel: channel,
  }
}

func (p *Process) SetOut(name string, channel chan interface{}) {
  p.outPorts[name] = &port{
    channel: channel,
  }
}

func (p *Process) Run() {
  go func(){
    task := p.task
    val := reflect.ValueOf(task).Elem()
    var wg sync.WaitGroup
    for name, p := range p.inPorts {
      wg.Add(1)
      go func(name string, p *port) {
        defer wg.Done()
        ch := p.channel
        v := reflect.ValueOf(<-ch)
        val.FieldByName(name).Set(v)
        if p.cache == nil{
          p.cache = v.Interface()
        }
        close(ch)
      }(name, p)
    }
    wg.Wait()

    // Execute the function of Process
    task.Execute()

    for name, port := range p.outPorts {
      wg.Add(1)
      go func(name string, ch chan interface{}) {
        defer wg.Done()
        v := val.FieldByName(name).Interface()
        ch <- v
        port.cache = v
        close(ch)
      }(name, port.channel)
    }
    wg.Wait()
  }()
}
