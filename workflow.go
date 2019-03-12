package flowmat

import (
  "log"
)

type port struct {
  channel chan interface{}
  cache interface{}
}

type Workflow struct {
  Name string
  proc map[string]*Process
  inPortss map[string]*port
  outPortss map[string]*port
}

func NewWorkflow(name string) *Workflow {
  wf := &Workflow{
    Name: name,
    proc: make(map[string]*Process),
    inPortss: make(map[string]*port),
    outPortss: make(map[string]*port),
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
  w.inPortss[name] = new(port)
  channel := make(chan interface{})
  port := w.inPortss[name]
  port.channel = channel

  p := w.proc[procName]
  p.inPorts[portName] = channel
}

func (w *Workflow) ExposeOut(name, procName, portName string) {
  w.outPortss[name] = new(port)
  channel := make(chan interface{})
  port := w.outPortss[name]
  port.channel = channel

  p := w.proc[procName]
  p.outPorts[portName] = channel
}

func (w *Workflow) In(portName string, data interface{}) {
  port := w.inPortss[portName]
  port.cache = data
}

func (w *Workflow) Out(portName string) interface{} {
  data := w.outPortss[portName].cache
  if data == nil {
    log.Panicf("%s has not get data", portName)
  }
  return data
}

func (w *Workflow) Run() {
  for _, p := range w.proc {
    p.Run()
  }
  for portName, port := range w.inPortss {
    cacheData := port.cache
    if cacheData == nil {
      log.Panicf("input not been set for port %s", portName)
    }
    port.channel <- cacheData
  }
  for _, port := range w.outPortss {
    data := <-port.channel
    port.cache = data
  }
}
