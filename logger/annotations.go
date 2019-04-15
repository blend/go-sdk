package logger

// Annotations are a collection of labels for an event.
type Annotations map[string]string

// AddAnnotationValue adds a label value.
func (a Annotations) AddAnnotationValue(key, value string) {
	a[key] = value
}

// GetAnnotationValue gets a label value.
func (a Annotations) GetAnnotationValue(key string) (value string, ok bool) {
	value, ok = a[key]
	return
}

// Decompose decomposes the labels into something we can write to json.
func (a Annotations) Decompose() map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range a {
		output[key] = value
	}
	return output
}
