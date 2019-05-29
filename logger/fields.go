package logger


// Fields are event meta fields.
type Fields = map[string]string

// CombineFields combines one or many set of fields.
func CombineFields(fields ...Fields) Fields {
	output := make(Fields)
	for _, fieldSet := range fields {
		for key, value := range fieldSet {
			output[key]=value
		}
	}
	return output
}