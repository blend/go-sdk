package jsonutil

// Flatten takes a map and returns an array of the values,
// discarding the map keys.
func Flatten(obj interface{}) (values []interface{}) {
	if typed, ok := obj.(map[string]interface{}); ok {
		for _, value := range typed {
			values = append(values, value)
		}
	}
	return
}
