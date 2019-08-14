package logger

// CombineLabels combines one or many set of fields.
func CombineLabels(labels ...Labels) Labels {
	output := make(Labels)
	for _, set := range labels {
		if set == nil || len(set) == 0 {
			continue
		}
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}

// Labels are a collection of labels for an event.
type Labels map[string]string

// Decompose decomposes the labels into something we can write to json.
func (l Labels) Decompose() map[string]interface{} {
	if l == nil {
		return nil
	}
	output := make(map[string]interface{})
	for key, value := range l {
		output[key] = value
	}
	return output
}
