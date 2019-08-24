package stats

// Taggable is an interface for specifying and retrieving default stats tags
type Taggable interface {
	AddDefaultTag(string, string)
	DefaultTags() []string
}
