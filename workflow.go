package flowmat

import (
  "log"
)

type Workflow struct {
  Name string
  proc map[string]*Process
  inPorts map[string]chan interface{}
  outPorts map[string]chan interface{}
  inCache map[string]interface{}
  outCache map[string]interface{}
}

func NewWorkflow(name string) *Workflow {
  wf := &Workflow{
    Name: name,
    proc: make(map[string]*Process),
    inPorts: make(map[string]chan interface{}),
    outPorts: make(map[string]chan interface{}),
    inCache: make(map[string]interface{}),
    outCache: make(map[string]interface{}),
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
  channel := make(chan interface{})
  w.inPorts[name] = channel

  p := w.proc[procName]
  p.inPorts[portName] = channel
}

func (w *Workflow) ExposeOut(name, procName, portName string) {
  channel := make(chan interface{})
  w.outPorts[name] = channel

  p := w.proc[procName]
  p.outPorts[portName] = channel
}

func (w *Workflow) In(portName string, data interface{}) {
  w.inCache[portName] = data
}

func (w *Workflow) Out(portName string) interface{} {
  data, ok := w.outCache[portName]
  if !ok {
    log.Panicf("%s has not get data", portName)
  }
  return data
}

func (w *Workflow) Run() {
  for portName, channel := range w.inPorts {
    cacheData, ok := w.inCache[portName]
    if !ok {
      log.Panicf("input not been set for port %s", portName)
    }
    channel <- cacheData
  }
  for _, p := range w.proc {
    p.Run()
  }
  for portName, channel := range w.outPorts {
    data := <-channel
    w.outCache[portName] = data
  }
}
