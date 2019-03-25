package flowmat

import (
  // "log"
  "fmt"
)

type Processer interface {
  Name() string
  Run()
  ExposeIn(string) *Port
  ExposeOut(string) *Port
}

type Workflow struct {
  name string
  proc map[string]Processer
  InPorts map[string]*Port
  OutPorts map[string]*Port
  exposePorts map[string]*Port
}

// NewWorkflow create workflow object
func NewWorkflow(name string) *Workflow {
  wf := &Workflow{
    name: name,
    proc: make(map[string]Processer),
    InPorts: make(map[string]*Port),
    OutPorts: make(map[string]*Port),
    exposePorts: make(map[string]*Port),
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

  out := s.ExposeOut(sendPort)
  in := r.ExposeIn(recvPort)

  go func() {
    v := <-out.channel
    in.channel <- v
  }()
}

// MapIn map inPorts of process to workflow
func (w *Workflow) MapIn(name, procName, portName string) {
  p := w.proc[procName]
  w.InPorts[name] = p.ExposeIn(portName)
}

// MapOut map outPorts of process to workflow
func (w *Workflow) MapOut(name, procName, portName string) {
  p := w.proc[procName]
  w.OutPorts[name] = p.ExposeOut(portName)
}

func (w *Workflow) ExposeIn(name string) *Port {
  w.exposePorts[name] = w.InPorts[name]
  return w.InPorts[name]
}

func (w *Workflow) ExposeOut(name string) *Port {
  w.exposePorts[name] = w.OutPorts[name]
  return w.OutPorts[name]
}

// In pass the data to the inport
func (w *Workflow) In(portName string, data interface{}) {
  p := w.InPorts[portName]
  p.cache = data
}

// // Out get the result from outport
// func (w *Workflow) Out(portName string) interface{} {
//   data := w.outPorts[portName].cache
//   // if data == nil {
//   //   log.Panicf("%s has not get data", portName)
//   // }
//   return data
// }

// Run the workflow aka its process in order
func (w *Workflow) Run() {
  for _, p := range w.proc {
    p.Run()
  }
  // for name, port := range w.InPorts {
  //   if _, ok := w.exposePorts[name]; !ok {
  //     port.channel <- port.cache
  //   }
  // }
  fmt.Println(w.InPorts["wfIn"])
  // w.InPorts["wfIn"].channel <- w.InPorts["wfIn"].cache
  // // if the port not expose, store it in cache
  // for name, port := range w.OutPorts {
  //   if _, ok := w.exposePorts[name]; !ok {
  //     data := <-port.channel
  //     port.cache = data
  //   }
  // }
}
