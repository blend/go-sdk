package r2

// OptMethod sets the request method.
func OptMethod(method string) Option {
	return func(r *Request) error {
		r.Method = method
		return nil
	}
}

// OptGet sets the request method.
func OptGet() Option {
	return func(r *Request) error {
		r.Method = "GET"
		return nil
	}
}

// OptPost sets the request method.
func OptPost() Option {
	return func(r *Request) error {
		r.Method = "POST"
		return nil
	}
}

// OptPut sets the request method.
func OptPut() Option {
	return func(r *Request) error {
		r.Method = "PUT"
		return nil
	}
}

// OptPatch sets the request method.
func OptPatch() Option {
	return func(r *Request) error {
		r.Method = "PATCH"
		return nil
	}
}

// OptDelete sets the request method.
func OptDelete() Option {
	return func(r *Request) error {
		r.Method = "DELETE"
		return nil
	}
}
