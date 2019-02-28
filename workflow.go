package giida

import (
  "reflect"
)

type Workflow struct {
  name string
  capacity int
  procs map[string]interface{}
  inPorts map[string]port
  outPorts map[string]port
  connections []connection
}

type port struct {
  proc string
  port string
  channel reflect.Value
}

type portName struct {
  proc string
  port string
}

type connection struct {
  src portName
  tgt portName
  channel reflect.Value
}

func NewWorkflow(name string, maxGoroutineTasks int) *Workflow {
  wf := &Workflow{
    name: name,
    capacity: maxGoroutineTasks,
    inPorts:  make(map[string]port, maxGoroutineTasks),
    outPorts: make(map[string]port, maxGoroutineTasks),
    connections:  make([]connection, 0, maxGoroutineTasks),
  }
  return wf
}
