package configutil

import "context"

var (
	_ StringsSource = (*StringsFunc)(nil)
)

// StringsFunc is a value source from a function.
type StringsFunc func(context.Context) ([]string, error)

// Strings returns an invocation of the function.
func (svf StringsFunc) Strings(ctx context.Context) ([]string, error) {
	return svf(ctx)
}
