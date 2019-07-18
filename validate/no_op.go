package validate

// NoOp is an empty validated.
type NoOp struct{}

// Validate implements the no op.
func (no NoOp) Validate() error { return nil }
