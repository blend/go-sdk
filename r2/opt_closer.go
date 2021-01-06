package r2

// OptCloser sets the request closer.
//
// It is typically used to clean up or trigger other actions.
func OptCloser(action func() error) Option {
	return func(r *Request) error {
		r.Closer = action
		return nil
	}
}
