package fsm

// InternalError is returned by FSM.Event() and should never occur. It is a
// probably because of a bug.
type InternalError struct{}

func (e InternalError) Error() string {
	return "internal error on state transition"
}

type InvalidEventError struct {
	Event string
	State string
}

func (e InvalidEventError) Error() string {
	return "event" + e.Event + " inappropriate in curruent state " + e.State
}

type UnknownEventError struct {
	Event string
}

func (e UnknownEventError) Error() string {
	return "event " + e.Event + " does not exist"
}
