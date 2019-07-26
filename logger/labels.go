package logger

// Labels are a collection of labels for an event.
type Labels map[string]string

// SetLabel adds a label value.
func (l Labels) SetLabel(key, value string) {
	l[key] = value
}

// GetLabelKeys returns the keys represented in the labels set.
func (l Labels) GetLabelKeys() (keys []string) {
	for key := range l {
		keys = append(keys, key)
	}
	return
}

// GetLabel gets a label value.
func (l Labels) GetLabel(key string) (value string, ok bool) {
	value, ok = l[key]
	return
}

// Decompose decomposes the labels into something we can write to json.
func (l Labels) Decompose() map[string]interface{} {
	output := make(map[string]interface{})
	for key, value := range l {
		output[key] = value
	}
	return output
}
