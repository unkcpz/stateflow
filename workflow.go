package stateflow

import (
	// "log"
	"fmt"
)

type Processer interface {
	Name() string
	Load()
	ExposeIn(string) *Port
	ExposeOut(string) *Port
}

type Workflow struct {
	name        string
	proc        map[string]Processer
	InPorts     map[string]*Port
	OutPorts    map[string]*Port
	exposePorts map[string]*Port
}

// NewWorkflow create workflow object
func NewWorkflow(name string) *Workflow {
	wf := &Workflow{
		name:        name,
		proc:        make(map[string]Processer),
		InPorts:     make(map[string]*Port),
		OutPorts:    make(map[string]*Port),
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

// NewProcess create and add process to workflow
func (w *Workflow) NewProcess(name string, task Tasker) *Process {
	p := NewProcess(name, task)
	w.Add(p)
	return p
}

// proc return proc of WF if not exist raise error
func (w *Workflow) getProc(name string) (p Processer, err error) {
	p, ok := w.proc[name]
	if !ok {
		return nil, fmt.Errorf("can't get Processer %s", name)
	}
	return p, nil
}

// Connect outport of Process A(sendProc) to inport of Process B(recvProc)
func (w *Workflow) Connect(sendProc, sendPort, recvProc, recvPort string) error {
	s, err := w.getProc(sendProc)
	if err != nil {
		return err
	}
	r, err := w.getProc(recvProc)
	if err != nil {
		return err
	}

	out := s.ExposeOut(sendPort)
	in := r.ExposeIn(recvPort)

	go func() {
		v := <-out.channel
		in.channel <- v
	}()

	return nil
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

// ExposeIn expose proc's port to workflow and return ptr to Port
func (w *Workflow) ExposeIn(name string) *Port {
	w.exposePorts[name] = w.InPorts[name]
	return w.InPorts[name]
}

// ExposeOut expose proc's port to workflow and return ptr to Port
func (w *Workflow) ExposeOut(name string) *Port {
	w.exposePorts[name] = w.OutPorts[name]
	return w.OutPorts[name]
}

// In pass the data to the inport
func (w *Workflow) In(portName string, data interface{}) {
	p := w.InPorts[portName]
	p.cache = data
}

// Out get the result from outport
func (w *Workflow) Out(portName string) interface{} {
	data := w.OutPorts[portName].cache
	return data
}

// Load the workflow aka its process in order, wait for input feeded in
func (w *Workflow) Load() {
	LogAuditf(w.Name(), "WF:Loading:[%s]", portInfo(w.InPorts))
	for _, p := range w.proc {
		p.Load()
	}
}

func (w *Workflow) Start() {
	LogAuditf(w.Name(), "WF:Running:[%s]", portInfo(w.InPorts))
	for name, port := range w.InPorts {
		if _, ok := w.exposePorts[name]; !ok {
			port.Feed(nil)
		}
	}
}

func (w *Workflow) Finish() {
	// if the port not expose, store it in cache
	for name, port := range w.OutPorts {
		if _, ok := w.exposePorts[name]; !ok {
			port.Extract()
		}
	}
	LogAuditf(w.Name(), "WF:Finished:[%s]", portInfo(w.OutPorts))
}

func (w *Workflow) Flow() {
	w.Load()
	w.Start()
	w.Finish()
}
