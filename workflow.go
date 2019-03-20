package flowmat

import (
  "log"
)

type Workflow struct {
  Name string
  proc map[string]*Process
  inPorts map[string]*port
  outPorts map[string]*port
}

func NewWorkflow(name string) *Workflow {
  wf := &Workflow{
    Name: name,
    proc: make(map[string]*Process),
    inPorts: make(map[string]*port),
    outPorts: make(map[string]*port),
  }
  return wf
}

func (w *Workflow) Add(p *Process) {
  w.proc[p.Name] = p
}

func (w *Workflow) Connect(sendProc, sendPort, recvProc, recvPort string) {
  s := w.proc[sendProc]
  r := w.proc[recvProc]
  out := make(chan interface{})
  in := make(chan interface{})

  s.SetOut(sendPort, out)
  r.SetIn(recvPort, in)

  go func() {
    v := <-out
    in <- v
  }()
}

func (w *Workflow) ExposeIn(name, procName, portName string) {
  w.inPorts[name] = new(port)
  channel := make(chan interface{})
  wfport := w.inPorts[name]
  wfport.channel = channel

  p := w.proc[procName]
  p.inPorts[portName] = &port{
    channel: channel,
  }
}

func (w *Workflow) ExposeOut(name, procName, portName string) {
  w.outPorts[name] = new(port)
  channel := make(chan interface{})
  wfport := w.outPorts[name]
  wfport.channel = channel

  p := w.proc[procName]
  p.outPorts[portName] = &port{
    channel: channel,
  }
}

func (w *Workflow) In(portName string, data interface{}) {
  port := w.inPorts[portName]
  port.cache = data
}

func (w *Workflow) Out(portName string) interface{} {
  data := w.outPorts[portName].cache
  if data == nil {
    log.Panicf("%s has not get data", portName)
  }
  return data
}

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
