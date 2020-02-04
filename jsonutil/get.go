package jsonutil

// Get follows a path of keys through an map which might have
// many generations of maps.
func Get(obj interface{}, path ...string) (interface{}, bool) {
	if len(path) == 0 {
		return nil, false
	}
	if typed, ok := obj.(map[string]interface{}); ok {
		return GetMap(typed, path...)
	}
	return nil, false
}

// GetMap follows a path of keys through an already type coerced map.
func GetMap(obj map[string]interface{}, path ...string) (interface{}, bool) {
	if len(path) == 0 {
		return nil, false
	}
	key := path[0]
	switch key {
	case FirstKey:
		if len(path) == 1 {
			values := Flatten(obj)
			if len(values) > 0 {
				return values[0], true
			}
			return nil, false
		}
		values := Flatten(obj)
		if len(values) > 0 {
			return Get(values[0], path[1:]...)
		}
		return nil, false
	case LastKey:
		if len(path) == 1 {
			values := Flatten(obj)
			if len(values) > 0 {
				return values[len(values)-1], true
			}
			return nil, false
		}
		values := Flatten(obj)
		if len(values) > 0 {
			return Get(len(values)-1, path[1:]...)
		}
		return nil, false
	default:
		if len(path) == 1 {
			value, ok := obj[key]
			return value, ok
		}
		value, ok := obj[key]
		if !ok {
			return nil, false
		}
		return Get(value, path[1:]...)
	}
}
