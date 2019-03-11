package giida

type Workflow struct {
  Name string
  proc map[string]*Process
}

func NewWorkflow(name string) {
  wf := &Workflow{
    Name: name,
  }
  return wf
}

func (w *Workflow) Add(p *Process) {
  w.proc[p.Name] = p
}

func (w *Workflow) Connect(sendProc, sendPort, recvProc, recvPort string) {
  // get value from send port
  s := w.proc[sendProc]
  ch := s.outPorts[sendPort]
  val := reflect.ValueOf(ch)


  r := w.proc[recvProc]
  r.SetIn(recvPort, channel)
}

func (w *Workflow) SetIn(procName, portName string, channel interface{}) {

}

func (w *Workflow) SetOut(procName, portName string, channel chan interface{})
