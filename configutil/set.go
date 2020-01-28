package configutil

import "time"

// SetString coalesces a given list of sources into a variable.
func SetString(destination *string, sources ...StringSource) error {
	var value *string
	var err error
	for _, source := range sources {
		value, err = source.String()
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

// SetStringPtr coalesces a given list of sources into a variable.
func SetStringPtr(destination **string, sources ...StringSource) error {
	var value *string
	var err error
	for _, source := range sources {
		value, err = source.String()
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

// SetStrings coalesces a given list of sources into a variable.
func SetStrings(destination *[]string, sources ...StringsSource) error {
	var value []string
	var err error
	for _, source := range sources {
		value, err = source.Strings()
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

// SetBool coalesces a given list of sources into a variable.
func SetBool(destination **bool, sources ...BoolSource) error {
	var value *bool
	var err error
	for _, source := range sources {
		value, err = source.Bool()
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

// SetInt coalesces a given list of sources into a variable.
func SetInt(destination *int, sources ...IntSource) error {
	var value *int
	var err error
	for _, source := range sources {
		value, err = source.Int()
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

// SetIntPtr coalesces a given list of sources into a variable.
func SetIntPtr(destination **int, sources ...IntSource) error {
	var value *int
	var err error
	for _, source := range sources {
		value, err = source.Int()
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

// SetFloat64 coalesces a given list of sources into a variable.
func SetFloat64(destination *float64, sources ...Float64Source) error {
	var value *float64
	var err error
	for _, source := range sources {
		value, err = source.Float64()
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

// SetFloat64Ptr coalesces a given list of sources into a variable.
func SetFloat64Ptr(destination **float64, sources ...Float64Source) error {
	var value *float64
	var err error
	for _, source := range sources {
		value, err = source.Float64()
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

// SetDuration coalesces a given list of sources into a variable.
func SetDuration(destination *time.Duration, sources ...DurationSource) error {
	var value *time.Duration
	var err error
	for _, source := range sources {
		value, err = source.Duration()
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

// SetDurationPtr coalesces a given list of sources into a variable.
func SetDurationPtr(destination **time.Duration, sources ...DurationSource) error {
	var value *time.Duration
	var err error
	for _, source := range sources {
		value, err = source.Duration()
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
