/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package configutil

import (
	"context"
	"time"
)

// SetString coalesces a given list of sources into a variable.
func SetString(destination *string, sources ...StringSource) ResolveAction {
	return func(ctx context.Context) error {
		var value *string
		var err error
		for _, source := range sources {
			value, err = source.String(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = *value
				return nil
			}
		}
		return nil
	}
}

// SetStringPtr coalesces a given list of sources into a variable.
func SetStringPtr(destination **string, sources ...StringSource) ResolveAction {
	return func(ctx context.Context) error {
		var value *string
		var err error
		for _, source := range sources {
			value, err = source.String(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = value
				return nil
			}
		}
		return nil
	}
}

// SetStrings coalesces a given list of sources into a variable.
func SetStrings(destination *[]string, sources ...StringsSource) ResolveAction {
	return func(ctx context.Context) error {
		var value []string
		var err error
		for _, source := range sources {
			value, err = source.Strings(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = value
				return nil
			}
		}
		return nil
	}
}

// SetBool coalesces a given list of sources into a variable.
func SetBool(destination *bool, sources ...BoolSource) ResolveAction {
	return func(ctx context.Context) error {
		var value *bool
		var err error
		for _, source := range sources {
			value, err = source.Bool(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = *value
				return nil
			}
		}
		return nil
	}
}

// SetBoolPtr coalesces a given list of sources into a variable.
func SetBoolPtr(destination **bool, sources ...BoolSource) ResolveAction {
	return func(ctx context.Context) error {
		var value *bool
		var err error
		for _, source := range sources {
			value, err = source.Bool(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = value
				return nil
			}
		}
		return nil
	}
}

// SetInt coalesces a given list of sources into a variable.
func SetInt(destination *int, sources ...IntSource) ResolveAction {
	return func(ctx context.Context) error {
		var value *int
		var err error
		for _, source := range sources {
			value, err = source.Int(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = *value
				return nil
			}
		}
		return nil
	}
}

// SetIntPtr coalesces a given list of sources into a variable.
func SetIntPtr(destination **int, sources ...IntSource) ResolveAction {
	return func(ctx context.Context) error {
		var value *int
		var err error
		for _, source := range sources {
			value, err = source.Int(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = value
				return nil
			}
		}
		return nil
	}
}

// SetInt32 coalesces a given list of sources into a variable.
func SetInt32(destination *int32, sources ...Int32Source) ResolveAction {
	return func(ctx context.Context) error {
		var value *int32
		var err error
		for _, source := range sources {
			value, err = source.Int32(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = *value
				return nil
			}
		}
		return nil
	}
}

// SetInt32Ptr coalesces a given list of sources into a variable.
func SetInt32Ptr(destination **int32, sources ...Int32Source) ResolveAction {
	return func(ctx context.Context) error {
		var value *int32
		var err error
		for _, source := range sources {
			value, err = source.Int32(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = value
				return nil
			}
		}
		return nil
	}
}

// SetInt64 coalesces a given list of sources into a variable.
func SetInt64(destination *int64, sources ...Int64Source) ResolveAction {
	return func(ctx context.Context) error {
		var value *int64
		var err error
		for _, source := range sources {
			value, err = source.Int64(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = *value
				return nil
			}
		}
		return nil
	}
}

// SetInt64Ptr coalesces a given list of sources into a variable.
func SetInt64Ptr(destination **int64, sources ...Int64Source) ResolveAction {
	return func(ctx context.Context) error {
		var value *int64
		var err error
		for _, source := range sources {
			value, err = source.Int64(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = value
				return nil
			}
		}
		return nil
	}
}

// SetFloat64 coalesces a given list of sources into a variable.
func SetFloat64(destination *float64, sources ...Float64Source) ResolveAction {
	return func(ctx context.Context) error {
		var value *float64
		var err error
		for _, source := range sources {
			value, err = source.Float64(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = *value
				return nil
			}
		}
		return nil
	}
}

// SetFloat64Ptr coalesces a given list of sources into a variable.
func SetFloat64Ptr(destination **float64, sources ...Float64Source) ResolveAction {
	return func(ctx context.Context) error {
		var value *float64
		var err error
		for _, source := range sources {
			value, err = source.Float64(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = value
				return nil
			}
		}
		return nil
	}
}

// SetDuration coalesces a given list of sources into a variable.
func SetDuration(destination *time.Duration, sources ...DurationSource) ResolveAction {
	return func(ctx context.Context) error {
		var value *time.Duration
		var err error
		for _, source := range sources {
			value, err = source.Duration(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = *value
				return nil
			}
		}
		return nil
	}
}

// SetDurationPtr coalesces a given list of sources into a variable.
func SetDurationPtr(destination **time.Duration, sources ...DurationSource) ResolveAction {
	return func(ctx context.Context) error {
		var value *time.Duration
		var err error
		for _, source := range sources {
			value, err = source.Duration(ctx)
			if err != nil {
				return err
			}
			if value != nil {
				*destination = value
				return nil
			}
		}
		return nil
	}
}
