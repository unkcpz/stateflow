package flowmat

import (
  "reflect"
  "sync"
  "fmt"
  "strings"
)

type Tasker interface{
  Execute()
}

type Process struct {
  name string
  task Tasker
  InPorts map[string]*Port
  OutPorts map[string]*Port
  ports map[string]*Port
  exposePorts map[string]*Port
  unsetPorts map[string]*Port
}

// NewProcess create a Process of a task
func NewProcess(name string, task Tasker) *Process {
  proc := &Process{
    name: name,
    task: task,
    InPorts: make(map[string]*Port),
    OutPorts: make(map[string]*Port),
    ports: make(map[string]*Port),
    exposePorts: make(map[string]*Port),
    unsetPorts: make(map[string]*Port),
  }
  val := reflect.ValueOf(task).Elem()
  // Bind every field of task to a port
  for i:=0; i<val.NumField(); i++ {
    fieldName := val.Type().Field(i).Name
    proc.ports[fieldName] = &Port{
      channel: make(chan interface{}),
    }
  }

  return proc
}

// Name return process's name
func (p  *Process) Name() string {
  return p.name
}

// ExposeIn return Port for channel
func (p *Process) ExposeIn(name string) *Port {
  port := p.ports[name]
  p.InPorts[name] = port
  p.exposePorts[name] = port
  return port
}

// ExposeOut return Port for channel
func (p *Process) ExposeOut(name string) *Port {
  port := p.ports[name]
  p.OutPorts[name] = port
  p.exposePorts[name] = port
  return port
}


// In set InPort's cache
func (p *Process) In(name string, data interface{}) {
  port := p.ports[name]
  port.cache = data
  p.InPorts[name] = port
}

// Out get OutPort's cache
func (p *Process) Out(name string) interface{} {
  port := p.ports[name]
  p.OutPorts[name] = port

  return port.cache
}

// collect unset ports
func collectUnsetPorts(proc *Process) {
  for name, _ := range proc.ports {
    _, okIn := proc.InPorts[name]
    _, okOut := proc.OutPorts[name]
    if !okIn && !okOut {
      proc.unsetPorts[name] = proc.ExposeOut(name)
    }
  }
}

func portInfo(ports map[string]*Port) string {
  var str strings.Builder
  var holder string
  for name, port := range ports {
    str.WriteString(holder)
    str.WriteString(name)
    str.WriteString(":")
    str.WriteString(fmt.Sprintf("%v", port.cache))
    holder = "; "
  }
  return str.String()
}

// Load process
func (p *Process) Load() {
  collectUnsetPorts(p)
  // gorountine get input and run Execute()
  go func(){
    task := p.task
    val := reflect.ValueOf(task).Elem()
    var wg sync.WaitGroup
    for name, port := range p.InPorts {
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

    LogAuditf(p.Name(), "PROC:Running:[%s]", portInfo(p.InPorts))

    // Execute the function of Process
    task.Execute()

    for name, port := range p.OutPorts {
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

    LogAuditf(p.Name(), "PROC:Finished:[%s]", portInfo(p.OutPorts))
  }()
}

// Start Process by feed all ready inputs
func (p *Process) Start() {
  // Feed the inputs aka start Chain
  for name, port := range p.InPorts {
    if _, ok := p.exposePorts[name]; !ok {
      port.Feed(nil)
    }
  }
}

// Finish Process by extract all output and store value to cache
func (p *Process) Finish() {
  // Extract the outputs
  for _, port := range p.unsetPorts {
    port.Extract()
  }
}
