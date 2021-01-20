package profanity

// RuleResult is a result from a rule.
type RuleResult struct {
	OK      bool
	File    string
	Line    int
	Message string
	Err     error
}
