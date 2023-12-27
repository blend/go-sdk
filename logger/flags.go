/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package logger

import (
	"strings"
)

// NewFlags returns a new flag set from an array of flag values.
// It applies some parsing rules, such as a `-` prefix denotes disabling the flag explicitly.
// `All` and `None` are special flag values that indicate all flags are enabled or none are enabled.
// Flags are caseless, and are lowercase in final output.
func NewFlags(flags ...string) *Flags {
	flagSet := &Flags{
		flags: make(map[string]bool),
	}

	for _, rawFlag := range flags {
		parsedFlag := strings.ToLower(strings.TrimSpace(rawFlag))
		if parsedFlag == FlagAll {
			flagSet.all = true
			continue
		}

		if parsedFlag == FlagNone {
			flagSet.none = true
			return flagSet
		}

		if strings.HasPrefix(parsedFlag, "-") {
			flagSet.flags[strings.TrimPrefix(parsedFlag, "-")] = false
		} else {
			flagSet.flags[parsedFlag] = true
		}
	}

	return flagSet
}

// FlagsAll returns a flags set with all enabled.
func FlagsAll() *Flags	{ return &Flags{all: true, flags: make(map[string]bool)} }

// FlagsNone returns a flags set with no flags enabled.
func FlagsNone() *Flags	{ return &Flags{none: true, flags: make(map[string]bool)} }

// Flags is a set of event flags.
type Flags struct {
	flags	map[string]bool
	all	bool
	none	bool
}

// Enable enables an event flag.
func (efs *Flags) Enable(flags ...string) {
	efs.none = false
	for _, flag := range flags {
		efs.flags[strings.ToLower(strings.TrimSpace(flag))] = true
	}
}

// Disable disables a flag.
func (efs *Flags) Disable(flags ...string) {
	for _, flag := range flags {
		efs.flags[strings.ToLower(strings.TrimSpace(flag))] = false
	}
}

// SetAll flips the `all` bit on the flag set to true.
// Note: flags that are explicitly disabled will remain disabled.
func (efs *Flags) SetAll() {
	efs.all = true
	efs.none = false
}

// All returns if the all bit is flipped to true.
func (efs *Flags) All() bool {
	return efs.all
}

// SetNone flips the `none` bit on the flag set to true.
// It also disables the `all` bit, and empties the enabled flag set.
func (efs *Flags) SetNone() {
	efs.all = false
	efs.flags = make(map[string]bool)
	efs.none = true
}

// None returns if the none bit is flipped to true.
func (efs *Flags) None() bool {
	return efs.none
}

// IsEnabled checks to see if an event is enabled.
func (efs Flags) IsEnabled(flag string) bool {
	switch {
	case efs.all:
		if efs.flags != nil {
			if enabled, hasEvent := efs.flags[flag]; hasEvent && !enabled {
				return false
			}
		}
		return true
	case efs.none:
		return false
	case efs.flags != nil:
		if enabled, hasFlag := efs.flags[flag]; hasFlag {
			return enabled
		}
	}
	return false
}

// String returns a string representation of the flags.
func (efs Flags) String() string {
	return strings.Join(efs.Flags(), ", ")
}

// Flags returns an array of flags.
func (efs Flags) Flags() []string {
	if efs.none {
		return []string{FlagNone}
	}

	var flags []string
	if efs.all {
		flags = []string{FlagAll}
	}
	for key, enabled := range efs.flags {
		if key != FlagAll {
			if enabled {
				if !efs.all {
					flags = append(flags, string(key))
				}
			} else {
				flags = append(flags, "-"+string(key))
			}
		}
	}
	return flags
}

// MergeWith sets the set from another, with the other taking precedence.
func (efs Flags) MergeWith(other *Flags) {
	if other.all {
		efs.all = true
	}
	if other.none {
		efs.none = true
	}
	for key, value := range other.flags {
		efs.flags[key] = value
	}
}
