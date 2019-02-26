package r2

import "io"

// OptBody sets the post body on the request.
func OptBody(contents io.ReadCloser) Option {
	return func(r *Request) error {
		r.Body = contents
		return nil
	}
}
