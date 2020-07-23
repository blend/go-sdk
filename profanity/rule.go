package profanity

import (
	"fmt"
	"strings"
)

// Rule is a serialized rule.
type Rule struct {
	// ID is a unique identifier for the rule.
	ID string `yaml:"id"`
	// File is the rules file path the rule came from.
	File string `yaml:"-"`
	// Description is a descriptive message for the rule.
	Description string `yaml:"description,omitempty"`

	// IncludeFiles sets a glob filter for file inclusion by filename.
	IncludeFiles []string `yaml:"includeFiles,omitempty"`
	// ExcludeFiles sets a glob filter for file exclusion by filename.
	ExcludeFiles []string `yaml:"excludeFiles,omitempty"`

	//
	// the below are matching rules.
	// if these match, the rule will fail the profanity check
	//

	// Contains implies we should fail if a file contains a given string.
	Contains []string `yaml:"contains,omitempty"`
	// Pattern implies we should fail if a file's content matches a given regex pattern.
	Pattern []string `yaml:"pattern,omitempty"`

	// go specific

	// ImportsContain enforces that a given list of imports are used.
	ImportsContain []string `yaml:"importsContain,omitempty"`
	// Calls enforces that a given list of imports are used.
	Calls []Call `yaml:"calls,omitempty"`
}

// ShouldInclude returns if we should include a file for a given rule.
// If the `.Include` field is unset, this will alway return true.
func (r Rule) ShouldInclude(file string) bool {
	if len(r.IncludeFiles) == 0 {
		return true
	}
	return GlobAnyMatch(r.IncludeFiles, file)
}

// ShouldExclude returns if we should include a file for a given rule.
// If the `.Include` field is unset, this will alway return true.
func (r Rule) ShouldExclude(file string) bool {
	// implicit rule:
	// we should omit non-go files from the imports ast parse
	if len(r.ImportsContain) > 0 {
		if !Glob(GoFiles, file) {
			return true
		}
	}
	if len(r.ExcludeFiles) == 0 {
		return false
	}
	return GlobAnyMatch(r.ExcludeFiles, file)
}

// Apply applies the rule.
func (r Rule) Apply(filename string, contents []byte) (result RuleResult) {
	if len(r.Contains) > 0 {
		result = ContainsAny(r.Contains...)(filename, contents)
		return
	}
	if len(r.Pattern) > 0 {
		result = MatchesAny(r.Pattern...)(filename, contents)
		return
	}
	if len(r.ImportsContain) > 0 {
		result = ImportsContainAny(r.ImportsContain...)(filename, contents)
		return
	}
	if len(r.Calls) > 0 {
		result = Calls(r.Calls...)(filename, contents)
		return
	}
	result = RuleResult{
		OK:   true,
		File: filename,
	}
	return
}

// String returns a string representation of the rule.
func (r Rule) String() string {
	var tokens []string

	if len(r.ID) > 0 {
		tokens = append(tokens, fmt.Sprintf("[%s]", r.ID))
	}
	if len(r.Description) > 0 {
		tokens = append(tokens, "`"+r.Description+"`")
	}
	if len(r.IncludeFiles) > 0 {
		tokens = append(tokens, fmt.Sprintf("[include files: %s]", strings.Join(r.IncludeFiles, ", ")))
	}
	if len(r.ExcludeFiles) > 0 {
		tokens = append(tokens, fmt.Sprintf("[exclude files: %s]", strings.Join(r.ExcludeFiles, ",")))
	}
	if len(r.Contains) > 0 {
		tokens = append(tokens, fmt.Sprintf("[contains: %s]", strings.Join(r.Contains, ",")))
	}
	if len(r.Pattern) > 0 {
		tokens = append(tokens, fmt.Sprintf("[matches patterns: %s]", strings.Join(r.Pattern, ",")))
	}
	if len(r.ImportsContain) > 0 {
		tokens = append(tokens, fmt.Sprintf("[go imports contain any: %s]", strings.Join(r.ImportsContain, ",")))
	}
	if len(r.Calls) > 0 {
		tokens = append(tokens, fmt.Sprintf("[go calls any: %s]", strings.Join(makeStrings(r.Calls), ",")))
	}
	return strings.Join(tokens, " ")
}

func makeStrings(objs ...interface{}) []string {
	var output []string
	for _, obj := range objs {
		output = append(output, fmt.Sprint(obj))
	}
	return output
}
