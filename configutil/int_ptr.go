package configutil

// IntPtr returns an IntSource for a given int pointer.
func IntPtr(value *int) IntSource {
	return IntPtrSource{Value: value}
}

var (
	_ IntSource = (*IntPtrSource)(nil)
)

// IntPtrSource is a IntSource that wraps an int pointer.
type IntPtrSource struct {
	Value *int
}

// Int implements IntSource.
func (ips IntPtrSource) Int() (*int, error) {
	return ips.Value, nil
}
