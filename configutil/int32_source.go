package configutil

import "context"

// Int32Source is a type that can return a value.
type Int32Source interface {
	// Int32 should return an int32 if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Int32(ctx context.Context) (*int32, error)
}

var (
	_ Int32Source = (*Int32)(nil)
)

// Int32 implements value provider.
//
// Note: Int32 treats 0 as unset, if 0 is a valid value you must use configutil.Int32Ptr.
type Int32 int32

// Int32 returns the value for a constant.
func (i Int32) Int32(_ context.Context) (*int32, error) {
	if i > 0 {
		value := int32(i)
		return &value, nil
	}
	return nil, nil
}
