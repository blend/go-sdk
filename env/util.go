package env

// PatchStringer is a type that handles unmarshalling a map of strings into itself.
type PatchStringer interface {
	PatchStrings(map[string]string) error
}
