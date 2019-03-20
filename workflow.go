package flowmat

import (
  "log"
)

type Processer interface {
  Name() string
  Run()
  SetIn(string, chan interface{})
  SetOut(string, chan interface{})
}

type Workflow struct {
  name string
  proc map[string]Processer
  inPorts map[string]*Port
  outPorts map[string]*Port
}

// NewWorkflow create workflow object
func NewWorkflow(name string) *Workflow {
  wf := &Workflow{
    name: name,
    proc: make(map[string]Processer),
    inPorts: make(map[string]*Port),
    outPorts: make(map[string]*Port),
  }
  return wf
}

func (w *Workflow) Name() string {
  return w.name
}

// Add process to workflow list
func (w *Workflow) Add(p Processer) {
  w.proc[p.Name()] = p
}

// Connect outport of Process A(sendProc) to inport of Process B(recvProc)
func (w *Workflow) Connect(sendProc, sendPort, recvProc, recvPort string) {
  s := w.proc[sendProc]
  r := w.proc[recvProc]
  out := make(chan interface{})
  in := make(chan interface{})

  // s.outPorts[sendPort] = &Port{
  //   channel: out,
  // }
  // r.inPorts[recvPort] = &Port{
  //   channel: in,
  // }
  s.SetOut(sendPort, out)
  r.SetIn(recvPort, in)

  go func() {
    v := <-out
    in <- v
  }()
}

// ExposeIn expose inPorts of process to workflow
func (w *Workflow) ExposeIn(name, procName, portName string) {
  w.inPorts[name] = new(Port)
  channel := make(chan interface{})
  wfport := w.inPorts[name]
  wfport.channel = channel

  p := w.proc[procName]
  // p.inPorts[portName] = &Port{
  //   channel: channel,
  // }
  p.SetIn(portName, channel)
}

// ExposeOut expose outPorts of process to workflow
func (w *Workflow) ExposeOut(name, procName, portName string) {
  w.outPorts[name] = new(Port)
  channel := make(chan interface{})
  wfport := w.outPorts[name]
  wfport.channel = channel

  p := w.proc[procName]
  // p.outPorts[portName] = &Port{
  //   channel: channel,
  // }
  p.SetOut(portName, channel)
}

// In pass the data to the inport
func (w *Workflow) In(portName string, data interface{}) {
  port := w.inPorts[portName]
  port.cache = data
}

// Out get the result from outport
func (w *Workflow) Out(portName string) interface{} {
  data := w.outPorts[portName].cache
  if data == nil {
    log.Panicf("%s has not get data", portName)
  }
  return data
}

// Run the workflow aka its process in order
func (w *Workflow) Run() {
  for _, p := range w.proc {
    p.Run()
  }
  for portName, port := range w.inPorts {
    cacheData := port.cache
    if cacheData == nil {
      log.Panicf("input not been set for port %s", portName)
    }
    port.channel <- cacheData
  }
  for _, port := range w.outPorts {
    data := <-port.channel
    port.cache = data
  }
}
