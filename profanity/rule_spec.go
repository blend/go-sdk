/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package profanity

import (
	"fmt"
	"strings"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/validate"
)

var (
	_ Rule = (*RuleSpec)(nil)
)

// Errors
const (
	ErrRuleSpecMultipleRules ex.Class = "rule spec invalid; multiple rule types specified"
	ErrRuleSpecRuleMissing   ex.Class = "rule spec invalid; at least one rule type is required"
)

// RuleSpec is a serialized rule.
type RuleSpec struct {
	// ID is a unique identifier for the rule.
	ID string `yaml:"id"`
	// SourceFile is the rules file path the rule came from.
	SourceFile string `yaml:"-"`
	// Description is a descriptive message for the rule.
	Description string `yaml:"description,omitempty"`
	// Files is the glob filter for inclusion and exclusion'
	// for this specific rule spec.
	Files GlobFilter `yaml:"files,omitempty"`
	// RuleSpecRules are the rules for the rule spec.
	RuleSpecRules `yaml:",inline"`
}

// Validate validates the RuleSpec.
func (r RuleSpec) Validate() error {
	if err := validate.String(&r.ID).Required()(); err != nil {
		return validate.Error(validate.ErrCause(err), r)
	}
	if err := r.Files.Validate(); err != nil {
		return validate.Error(validate.ErrCause(err), r)
	}
	if err := r.RuleSpecRules.Validate(); err != nil {
		return validate.Error(validate.ErrCause(err), r)
	}
	return nil

}

// String returns a string representation of the rule.
func (r RuleSpec) String() string {
	var tokens []string

	if len(r.ID) > 0 {
		tokens = append(tokens, r.ID)
	}
	if len(r.Description) > 0 {
		tokens = append(tokens, "`"+r.Description+"`")
	}
	tokens = append(tokens, r.Files.String())
	tokens = append(tokens, r.RuleSpecRules.String())
	return fmt.Sprintf("[%s]", strings.Join(tokens, " "))
}

// RuleSpecRules are the specific rules for a given RuleSpec.
//
// The usage of this should be that only _one_ of these rules
// should be set for a given rule spec.
type RuleSpecRules struct {
	// Contains implies we should fail if a file contains a given string.
	Contents *Contents `yaml:"contents,omitempty"`
	// GoImportsContain enforces that a given list of imports are used.
	GoImports *GoImports `yaml:"goImports,omitempty"`
	// GoCalls enforces that a given list of imports are used.
	GoCalls *GoCalls `yaml:"goCalls,omitempty"`
	// NoGenericDecls enforces that no generic declarations exist.
	NoGenericDecls *NoGenericDecls `yaml:"noGenericDecls,omitempty"`
}

// Rules returns the rules from the spec.
//
// Note: you should add new rule types here and on the type itself.
func (r RuleSpecRules) Rules() (output []Rule) {
	if r.Contents != nil {
		output = append(output, r.Contents)
	}
	if r.GoImports != nil {
		output = append(output, r.GoImports)
	}
	if r.GoCalls != nil {
		output = append(output, r.GoCalls)
	}
	if r.NoGenericDecls != nil {
		output = append(output, r.NoGenericDecls)
	}
	return
}

// Rule returns the active rule from the spec.
func (r RuleSpecRules) Rule() Rule {
	if rules := r.Rules(); len(rules) > 0 {
		return rules[0]
	}
	return nil
}

// Check applies the rule.
func (r RuleSpecRules) Check(filename string, contents []byte) (result RuleResult) {
	if result = r.Rule().Check(filename, contents); !result.OK {
		return
	}
	result = RuleResult{
		OK:   true,
		File: filename,
	}
	return
}

// Validate validates the rule spec rules.
func (r RuleSpecRules) Validate() error {
	if len(r.Rules()) > 1 {
		return validate.Error(ErrRuleSpecMultipleRules, nil)
	}
	if len(r.Rules()) == 0 {
		return validate.Error(ErrRuleSpecRuleMissing, nil)
	}
	if typed, ok := r.Rule().(interface {
		Validate() error
	}); ok {
		if err := typed.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// String implements fmt.Stringer.
func (r RuleSpecRules) String() string {
	var tokens []string
	for _, rule := range r.Rules() {
		if typed, ok := rule.(fmt.Stringer); ok {
			tokens = append(tokens, typed.String())
		}
	}
	return strings.Join(tokens, " ")
}
