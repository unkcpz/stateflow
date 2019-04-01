package fsm

// InternalError is returned by FSM.Event() and should never occur. It is a
// probably because of a bug.
type InternalError struct{}

func (e InternalError) Error() string {
	return "internal error on state transition"
}
