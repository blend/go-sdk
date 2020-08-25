package configutil

import "context"

// Int64Source is a type that can return a value.
type Int64Source interface {
	// Int should return a int if the source has a given value.
	// It should return nil if the value is not found.
	// It should return an error if there was a problem fetching the value.
	Int64(ctx context.Context) (*int64, error)
}

var (
	_ Int64Source = (*Int64)(nil)
)

// Int64 implements value provider.
//
// Note: Int64 treats 0 as unset, if 0 is a valid value you must use configutil.Int64Ptr.
type Int64 int64

// Int64 returns the value for a constant.
func (i Int64) Int64(_ context.Context) (*int64, error) {
	if i > 0 {
		value := int64(i)
		return &value, nil
	}
	return nil, nil
}
