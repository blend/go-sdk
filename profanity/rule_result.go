package profanity

import (
	"fmt"
	"strings"
	"github.com/blend/go-sdk/ansi"
)

// RuleResult is a result from a rule.
type RuleResult struct {
	OK      bool
	File    string
	Line    int
	Message string
	Err     error
}

// Failure returns a failure error message for a given file and error.
func (r RuleResult) Failure(rule Rule) error {
	var tokens []string
	tokens = append(tokens, fmt.Sprintf("%s:%d", ansi.Bold(ansi.ColorWhite, r.File), r.Line))
	if rule.ID != "" {
		tokens = append(tokens, fmt.Sprintf("\t%s: %s", ansi.LightBlack("id"), rule.ID))
	}
	if rule.Description != "" {
		tokens = append(tokens, fmt.Sprintf("\t%s: %s", ansi.LightBlack("description"), rule.Description))
	}
	tokens = append(tokens, fmt.Sprintf("\t%s: %s", ansi.LightBlack("status"), ansi.Red("failed")))
	tokens = append(tokens, fmt.Sprintf("\t%s: %s", ansi.LightBlack("rule"), r.Message))
	return fmt.Errorf(strings.Join(tokens, "\n"))
}
