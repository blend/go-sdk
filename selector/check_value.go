package selector

// CheckValue returns if the value is valid.
func CheckValue(value string) error {
	if len(value) > MaxLabelValueLen {
		return ErrLabelValueTooLong
	}
	return checkName(value)
}
