package configutil

import "context"

// IntSource is a type that can return a value.
type IntSource interface {
	// Int should return a int if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Int(ctx context.Context) (*int, error)
}

var (
	_ IntSource = (*Int)(nil)
)

// Int implements value provider.
//
// Note: Int treats 0 as unset, if 0 is a valid value you must use configutil.IntPtr.
type Int int

// Int returns the value for a constant.
func (i Int) Int(_ context.Context) (*int, error) {
	if i > 0 {
		value := int(i)
		return &value, nil
	}
	return nil, nil
}
