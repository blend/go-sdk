package selector

// Labels is an alias for map[string]string
type Labels = map[string]string

// Selector is the common interface for selector types.
type Selector interface {
	Matches(labels Labels) bool
	Validate() error
	String() string
}
