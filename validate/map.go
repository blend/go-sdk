package validate

import (
	"reflect"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/reflectutil"
)

// Errors
const (
	ErrInstanceNotMap ex.Class = "validated reference is not a map"

	ErrMapKeys ex.Class = "map should have keys"
)

// Map returns map validators.
func Map(instance interface{}) MapValidators {
	return MapValidators{instance}
}

// MapValidators is a set of validators for maps.
type MapValidators struct {
	Instance interface{}
}

// Keys validates a map contains a given set of keys.
func (mv MapValidators) Keys(keys ...interface{}) Validator {
	return func() error {
		value := reflectutil.Value(mv.Instance)
		if value.Kind() != reflect.Map {
			return ErrInstanceNotMap
		}

		for _, key := range keys {
			mapValue := value.MapIndex(reflect.ValueOf(key))
			if !mapValue.IsValid() {
				return Errorf(ErrMapKeys, "missing key: %v", key)
			}
		}
		return nil
	}
}
