package giida

import (
  "errors"
  "fmt"
  "reflect"
  "sync"
)

type Workflow struct {
  name string
  capacity int
  bufferSize int
  waitGrp *sync.WaitGroup
  procs map[string]Process
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

func NewWorkflow(name string, maxGoroutineTasks, bufferSize int) *Workflow {
  wf := &Workflow{
    name: name,
    capacity: maxGoroutineTasks,
    bufferSize: bufferSize,
    waitGrp:  new(sync.WaitGroup),
    procs:  make(map[string]Process, maxGoroutineTasks),
    inPorts:  make(map[string]port, maxGoroutineTasks),
    outPorts: make(map[string]port, maxGoroutineTasks),
    connections:  make([]connection, 0, maxGoroutineTasks),
  }
  return wf
}

// Add adds a new process with a given name to the network.
func (n *Workflow) Add(name string, proc Process) error {
  n.procs[name] = proc
  return nil
}

func (n *Workflow) Connect(senderName, senderPort, receiverName, receiverPort string) error {
  return n.ConnectBuf(senderName, senderPort, receiverName, receiverPort, n.bufferSize)
}

// Connect connects a sender to a receiver and creates a channel between them
func (n *Workflow) ConnectBuf(senderName, senderPort, receiverName, receiverPort string, bufferSize int) error {
  senderPortVal, err := n.getProcPort(senderName, senderPort, reflect.SendDir)
  if err != nil {
    return err
  }

  receiverPortVal, err := n.getProcPort(receiverName, receiverPort, reflect.RecvDir)
  if err != nil {
    return err
  }

  var channel reflect.Value
  // if !receiverPortVal.IsNil() {
  //   // Find existing channel attached to the receiver
  //   channel =
  // }

  sndPortType := senderPortVal.Type()
  // Create a new channel if none of the existing channels found
  if !channel.IsValid() {
    // Make a channel of an appropriate type
    chanType := reflect.ChanOf(reflect.BothDir, sndPortType.Elem())
    channel = reflect.MakeChan(chanType, bufferSize)
  }
  // Set the channels
  if senderPortVal.IsNil() {
    senderPortVal.Set(channel)
  }
  if receiverPortVal.IsNil() {
    receiverPortVal.Set(channel)
  }

  // Add connection info
  n.connections = append(n.connections, connection{
    src: portName{proc: senderName, port: senderPort},
    tgt: portName{proc: receiverName, port: receiverPort},
    channel: channel,
  })
  return nil
}

func (n *Workflow) getProcPort(procName, portName string, dir reflect.ChanDir) (reflect.Value, error) {
  nilValue := reflect.ValueOf(nil)
  // Ensure process exists
  proc, ok := n.procs[procName]
  if !ok {
    return nilValue, fmt.Errorf("Connect error: process '%s' not found", procName)
  }

  // Ensure sender is settable
  val := reflect.ValueOf(proc)
  if val.Kind() == reflect.Ptr && val.IsValid() {
    val = val.Elem()
  }
  if !val.CanSet() {
    return nilValue, fmt.Errorf("Connect error: process '%s' is not settable", procName)
  }

  var portVal reflect.Value
  var err error
  portVal = val.FieldByName(portName)
  if !portVal.IsValid() {
    err = errors.New("")
  }
  if err != nil {
    return nilValue, fmt.Errorf("Connect error: process '%s' does not have port '%s'", procName, portName)
  }

  return portVal, nil
}

// Task runs the net
func (n *Workflow) Task() {
  for _, p := range n.procs {
    n.waitGrp.Add(1)
    wait := Run(p)
    go func() {
      <-wait
      n.closeProcOuts(p)
      n.waitGrp.Done()
    }()
  }
  n.waitGrp.Wait()
}

func (n *Workflow) closeProcOuts(proc Process) {
  val := reflect.ValueOf(proc).Elem()
  for i := 0; i < val.NumField(); i++ {
    field := val.Field(i)
    fieldType := field.Type()
    if !(field.IsValid() && field.Kind() == reflect.Chan && field.CanSet() &&
      fieldType.ChanDir()&reflect.SendDir != 0 && fieldType.ChanDir()&reflect.RecvDir == 0) {
        continue
    }
  }
}

func (n *Workflow) getInPort(name string) (reflect.Value, error) {
  if pName, ok := n.inPorts[name]; ok {
    return pName.channel, nil
  }
  return reflect.ValueOf(nil), fmt.Errorf("Inport not found: '%s'", name)
}

func (n *Workflow) getOutPort(name string) (reflect.Value, error) {
  if pName, ok := n.outPorts[name]; ok {
    return pName.channel, nil
  }
  return reflect.ValueOf(nil), fmt.Errorf("Outport not found: '%s'", name)
}

// MapInPort adds an inport to the net and maps it to a contained proc's port
func (n *Workflow) MapInPort(name, procName, procPort string) error {
  var channel reflect.Value
  var err error
  if _, found := n.procs[procName]; !found {
    return fmt.Errorf("Could not map inport: process '%s' not found", procName)
  }
  channel, err = n.getProcPort(procName, procPort, reflect.RecvDir)
  if err != nil {
    return err
  }
  n.inPorts[name] = port{proc: procName, port: procPort, channel: channel}
  return nil
}

func (n *Workflow) MapOutPort(name, procName, procPort string) error {
  var channel reflect.Value
  var err error
  if _, found := n.procs[procName]; !found {
    return fmt.Errorf("Could not map outport: process '%s' not found", procName)
  }
  channel, err = n.getProcPort(procName, procPort, reflect.SendDir)
  if err != nil {
    return err
  }
  n.outPorts[name] = port{proc: procName, port: procPort, channel: channel}
  return nil
}

// SetInPort assigns a channel to a network's inport to talk to the outer world
func (n *Workflow) SetInPort(name string, channel interface{}) error {
  p, err := n.getInPort(name)
  if err != nil {
    return err
  }
  // Try to set it
  if p.CanSet() {
    p.Set(reflect.ValueOf(channel))
  } else {
    return fmt.Errorf("Cannot set graph inport: '%s'", name)
  }

  // Save it in inPorts to be used with IIPs if needed
  if p, ok := n.inPorts[name]; ok {
    p.channel = reflect.ValueOf(channel)
    n.inPorts[name] = p
  }
  return nil
}

// SetOutPort assigns a channel to a network's outport to talk to the outer world
func (n *Workflow) SetOutPort(name string, channel interface{}) error {
  p, err := n.getOutPort(name)
  if err != nil {
    return err
  }
  // Try to set it
  if p.CanSet() {
    p.Set(reflect.ValueOf(channel))
  } else {
    return fmt.Errorf("Cannot set graph inport: '%s'", name)
  }

  // Save it in outPorts to be used with IIPs if needed
  if p, ok := n.outPorts[name]; ok {
    p.channel = reflect.ValueOf(channel)
    n.outPorts[name] = p
  }
  return nil
}
