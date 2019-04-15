package logger

// Labels are a collection of labels for an event.
type Labels map[string]string

// AddLabelValue adds a label value.
func (l Labels) AddLabelValue(key, value string) {
	l[key] = value
}

// GetLabelValue gets a label value.
func (l Labels) GetLabelValue(key string) (value string, ok bool) {
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
