package fsm

type transitioner interface {
  transition(*FSM) error
}

type FSM struct {
  current string
  transitions map[eKey]string
  transition func()
  transitionerObj transitioner
}

type EventDesc struct {
  Name string
  Src []string
  Dst string
}

type Events []EventDesc

type Callback func(*Event)

type Callbacks map[string]Callback

func NewFSM(initial string, events []EventDesc, callbacks map[string]Callback) *FSM {
  f := &FSM {
    transitionerObj: &transitionerStruct{},
    current: initial,
    transitions: make(map[eKey]string),
  }
  // Build transition map and store sets of all events and states
  allEvents := make(map[string]bool)
  allStates := make(map[string]bool)
  for _, e := range events {
    for _, src := range e.Src {
      f.transitions[eKey{e.Name, src}] = e.Dst
      allStates[src] = true
      allStates[e.Dst] = true
    }
    allEvents[e.Name] = true
  }
  return f
}

func (f *FSM) Event(event string) error {
  // f.eventMu.Lock()
  // defer f.evetMu.Unlock()
  //
  // f.stateMu.RLock()
  // defer f.stateMu.RUnlock()
  dst, ok := f.transitions[eKey{event, f.current}]
  if !ok {
    for ekey := range f.transitions {
      if ekey.event == event {
        return InvalidEventError{event, f.current}
      }
    }
    return UnknownEventError{event}
  }
  f.transition = func() {
    f.current = dst
  }

  err := f.doTransition()
  if err != nil {
    return InternalError{}
  }
  return nil
}

func (f *FSM) doTransition() error {
  return f.transitionerObj.transition(f)
}

func (f *FSM) Current() string {
  // f.stateMu.RLock()
  // defer f.stateMu.RUnlock()
  return f.current
}

func (f *FSM) SetState(state string) {
  f.current = state
  return
}

type transitionerStruct struct{}

func (t transitionerStruct) transition(f *FSM) error {
  f.transition()
  return nil
}

const (
  callbackNone int = iota
  callbackBeforeEvent
  callbackLeaveState
  callbackEnterState
  callbackAfterEvent
)

type cKey struct  {
  target string
  callbackType int
}

type eKey struct {
  event string
  src string
}
