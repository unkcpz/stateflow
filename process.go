package flowmat

import (
  "reflect"
  "sync"
  // "fmt"
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

func (p *Process) In(name string, data interface{}) {
  p.inPorts[name] = &Port{
    channel: make(chan interface{}),
    cache: data,
  }
}

func (p *Process) Out(name string) interface{} {
  port := p.outPorts[name]
  return port.cache
}

func (p *Process) Run() {
  task := p.task
  val := reflect.ValueOf(task).Elem()
  // unset is Field of unset Ports
  unset := make([]string, 0)
  for i:=0; i<val.NumField(); i++ {
    fieldName := val.Type().Field(i).Name
    // set only when fieldName not set in inPorts or outPorts
    _, okIn := p.inPorts[fieldName]
    _, okOut := p.outPorts[fieldName]
    if !okIn && !okOut {
      p.outPorts[fieldName] = &Port{
        channel: make(chan interface{}),
      }
      unset = append(unset, fieldName)
    }
  }

  // gorountine get input and run Execute()
  go func(){
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
        port.cache = v
        ch <- v
        close(ch)
      }(name, port)
    }
    wg.Wait()
  }()

  // 
  for _, port := range p.inPorts {
    cacheData := port.cache
    if cacheData != nil {
      port.channel <- cacheData
    }
  }

  for _, name := range unset {
    // cache is already set in Run
    <-p.outPorts[name].channel
  }
}
