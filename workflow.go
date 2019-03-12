package giida

import (
  // "reflect"
  // "fmt"
)

type Workflow struct {
  Name string
  proc map[string]*Process
}

func NewWorkflow(name string) *Workflow {
  wf := &Workflow{
    Name: name,
    proc: make(map[string]*Process),
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
  // val := reflect.ValueOf(r.task).Elem()
  // fmt.Println("!!!")
  // fmt.Println(reflect.ValueOf(val.FieldByName(recvPort)))

  // fmt.Println("!!!")
  // chanType := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(val.FieldByName(recvPort)))
  // chanType := reflect.ChanOf(reflect.BothDir, reflect.TypeOf(1))
  // in := reflect.MakeChan(chanType, 0)
  in := make(chan interface{})

  s.SetOut(sendPort, out)
  r.SetIn(recvPort, in)

  go func() {
    v := <-out
    // <- out
    // fmt.Println(<-out)
    // in.Send(reflect.ValueOf(1))
    in <- v
  }()
}

func (w *Workflow) SetIn(procName, portName string, channel chan interface{}) {
  p := w.proc[procName]
  p.inPorts[portName] = channel
}

func (w *Workflow) SetOut(procName, portName string, channel chan interface{}) {
  p := w.proc[procName]
  p.outPorts[portName] = channel
}

func (w *Workflow) Run() {
  for _, p := range w.proc {
    p.Run()
  }
}
