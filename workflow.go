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
  ips reflect.Value
}

type portName struct {
  proc string
  port string
}

type connection struct {
  src portName
  tgt portName
  ips reflect.Value
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

// Add adds a new Task with a given name to the network.
func (n *Workflow) Add(proc *Process) error {
	name := proc.Name
	n.procs[name] = *proc
  return nil
}

// func (n *Workflow) Connect(senderName, senderPort, receiverName, receiverPort string) error {
//   return n.ConnectBuf(senderName, senderPort, receiverName, receiverPort, n.bufferSize)
// }

// Connect connects a sender to a receiver and creates a channel between them
func (n *Workflow) Connect(senderName, senderPort, receiverName, receiverPort string) error {
  senderPortVal, err := n.getProcPort(senderName, senderPort, reflect.SendDir)
  if err != nil {
    return err
  }

  receiverPortVal, err := n.getProcPort(receiverName, receiverPort, reflect.RecvDir)
  if err != nil {
    return err
  }

  var ips reflect.Value

  sndPortType := senderPortVal.Type()
  // Create a new channel if none of the existing channels found
  if !ips.IsValid() {
    // Make a channel of an appropriate type
    chanType := reflect.ChanOf(reflect.BothDir, sndPortType.Elem())
    ips = reflect.MakeChan(chanType, 2)
  }
  // Set the channels
  if senderPortVal.IsNil() {
    senderPortVal.Set(ips)
  }
  if receiverPortVal.IsNil() {
    receiverPortVal.Set(ips)
  }

  // Add connection info
  n.connections = append(n.connections, connection{
    src: portName{proc: senderName, port: senderPort},
    tgt: portName{proc: receiverName, port: receiverPort},
    ips: ips,
  })
  return nil
}

func (n *Workflow) getProcPort(procName, portName string, dir reflect.ChanDir) (reflect.Value, error) {
  nilValue := reflect.ValueOf(nil)
  // Ensure Process exists
  proc, ok := n.procs[procName]
  task := proc.task
  if !ok {
    return nilValue, fmt.Errorf("Connect error: Task '%s' not found", procName)
  }

  // Ensure sender is settable
  val := reflect.ValueOf(task)
  if val.Kind() == reflect.Ptr && val.IsValid() {
    val = val.Elem()
  }
  if !val.CanSet() {
    return nilValue, fmt.Errorf("Connect error: Task '%s' is not settable", procName)
  }

  var portVal reflect.Value
  var err error
  portVal = val.FieldByName(portName)
  if !portVal.IsValid() {
    err = errors.New("")
  }
  if err != nil {
    return nilValue, fmt.Errorf("Connect error: Task '%s' does not have port '%s'", procName, portName)
  }

  return portVal, nil
}

// Execute runs the net
func (n *Workflow) Execute() {
  for _, p := range n.procs {
    n.waitGrp.Add(1)
    wait := p.Run()
    go func() {
      defer n.waitGrp.Done()
      <-wait
			// p channel already closed
      // n.closeProcOuts(p)
    }()
  }
  n.waitGrp.Wait()
}

func (n *Workflow) Run() Wait {
	wait := make(Wait)
	go func() {
		// fmt.Printf("%s | Running %s\n", timeStamp(), p.Name)
		n.Execute()

		wait <- Done{}
		// fmt.Printf("%s | %s Finished\n", timeStamp(), p.Name)
	}()
	return wait
}

// func (n *Workflow) closeProcOuts(proc Process) {
//   val := reflect.ValueOf(proc.task).Elem()
//   for i := 0; i < val.NumField(); i++ {
//     field := val.Field(i)
//     fieldType := field.Type()
//     if !(field.IsValid() && field.Kind() == reflect.Chan && field.CanSet() &&
//       fieldType.ChanDir()&reflect.SendDir != 0 && fieldType.ChanDir()&reflect.RecvDir == 0) {
//         continue
//     }
//     // field.Close()
//   }
// }

func (n *Workflow) getInPort(name string) (reflect.Value, error) {
  if pName, ok := n.inPorts[name]; ok {
    return pName.ips, nil
  }
  return reflect.ValueOf(nil), fmt.Errorf("Inport not found: '%s'", name)
}

func (n *Workflow) getOutPort(name string) (reflect.Value, error) {
  if pName, ok := n.outPorts[name]; ok {
    return pName.ips, nil
  }
  return reflect.ValueOf(nil), fmt.Errorf("Outport not found: '%s'", name)
}

// MapInPort adds an inport to the net and maps it to a contained proc's port
func (n *Workflow) MapInPort(name, procName, procPort string) error {
  var ips reflect.Value
  var err error
  if _, found := n.procs[procName]; !found {
    return fmt.Errorf("Could not map inport: Proc '%s' not found", procName)
  }
  ips, err = n.getProcPort(procName, procPort, reflect.RecvDir)
  if err != nil {
    return err
  }
  n.inPorts[name] = port{proc: procName, port: procPort, ips: ips}
  return nil
}

func (n *Workflow) MapOutPort(name, procName, procPort string) error {
  var ips reflect.Value
  var err error
  if _, found := n.procs[procName]; !found {
    return fmt.Errorf("Could not map outport: Proc '%s' not found", procName)
  }
  ips, err = n.getProcPort(procName, procPort, reflect.SendDir)
  if err != nil {
    return err
  }
  n.outPorts[name] = port{proc: procName, port: procPort, ips: ips}
  return nil
}

// SetInPort assigns a channel to a network's inport to talk to the outer world
func (n *Workflow) SetInPort(name string, ips interface{}) error {
  p, err := n.getInPort(name)
  if err != nil {
    return err
  }
  // Try to set it
  if p.CanSet() {
    p.Set(reflect.ValueOf(ips))
  } else {
    return fmt.Errorf("Cannot set graph inport: '%s'", name)
  }

  // Save it in inPorts to be used with IIPs if needed
  if p, ok := n.inPorts[name]; ok {
    p.ips = reflect.ValueOf(ips)
    n.inPorts[name] = p
  }
  return nil
}

// SetOutPort assigns a channel to a network's outport to talk to the outer world
func (n *Workflow) SetOutPort(name string, ips interface{}) error {
  p, err := n.getOutPort(name)
  if err != nil {
    return err
  }
  // Try to set it
  if p.CanSet() {
    p.Set(reflect.ValueOf(ips))
  } else {
    return fmt.Errorf("Cannot set graph inport: '%s'", name)
  }

  // Save it in outPorts to be used with IIPs if needed
  if p, ok := n.outPorts[name]; ok {
    p.ips = reflect.ValueOf(ips)
    n.outPorts[name] = p
  }
  return nil
}
