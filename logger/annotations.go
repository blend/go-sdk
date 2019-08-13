package logger

// CombineAnnotations combines one or many set of annotations.
func CombineAnnotations(annotations ...Annotations) Annotations {
	output := make(Annotations)
	for _, set := range annotations {
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}

// Annotations are a collection of labels for an event.
type Annotations map[string]interface{}

// AddAnnotationValue adds a label value.
func (a Annotations) AddAnnotationValue(key string, value interface{}) {
	a[key] = value
}

// GetAnnotationKeys returns the keys represented in the annotations set.
func (a Annotations) GetAnnotationKeys() (keys []string) {
	for key := range a {
		keys = append(keys, key)
	}
	return
}

// GetAnnotationValue gets a label value.
func (a Annotations) GetAnnotationValue(key string) (value interface{}, ok bool) {
	value, ok = a[key]
	return
}

// Decompose decomposes the labels into something we can write to json.
func (a Annotations) Decompose() map[string]interface{} {
	if a == nil {
		return nil
	}
	output := make(map[string]interface{})
	for key, value := range a {
		output[key] = value
	}
	return output
}
