package configutil

import "context"

// StringSource is a type that can return a value.
type StringSource interface {
	// String should return a string if the source has a given value.
	// It should return nil if the value is not present.
	// It should return an error if there was a problem fetching the value.
	String(context.Context) (*string, error)
}

var (
	_ StringSource = (*String)(nil)
)

// String implements value provider.
// An empty string is treated as unset, and will cause `String(context.Context) (*string, error)` to return
// nil for the string.
type String string

// StringValue returns the value for a constant.
func (s String) String(_ context.Context) (*string, error) {
	value := string(s)
	if value == "" {
		return nil, nil
	}
	return &value, nil
}
