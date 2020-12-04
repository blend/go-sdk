package profanity

// Rule is a criteria for profanity.
type Rule interface {
	Check(file string, contents []byte) RuleResult
}
