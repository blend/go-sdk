package logger

// CombineFields combines one or many set of fields.
func CombineFields(fields ...Fields) Fields {
	output := make(Fields)
	for _, set := range fields {
		if set == nil || len(set) == 0 {
			continue
		}
		for key, value := range set {
			output[key] = value
		}
	}
	return output
}

// Fields are a collection of extra context fields for an event.
type Fields map[string]interface{}
