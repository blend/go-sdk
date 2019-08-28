package env

import (
	"fmt"
	"unicode"

	"strings"

	"bytes"
	"github.com/blend/go-sdk/ex"
)

const (
	// ValueDelimiter ("=") is the delimiter between a key and a value for an
	// environment variable.
	valueDelimiter = "="

	// QuoteDelimiter (`"`) is a delimiter indicating a string literal. This
	// gives the user the option to have spaces, for example, in their
	// environment variable values.
	quoteDelimiter = "\""

	// EscapeDelimiter ("\") is used to escape the next character so it is
	// accepted as a part of the input value.
	escapeDelimiter = "\\"
)

// DFAState is a wrapper type for the standard enum integer type, representing
// the state of the parsing table for the DFA. We create a new type so that we
// can use a switch case on this particular enum type and not worry about
// accidentally setting the state to an invalid value.
type dfaState int

const (
	// rootState is the "default" starting state state. It processes text
	// normally, performing actions on tokens and excluding whitespace.
	rootState dfaState = iota

	// escapeState represents the state encountered after the parser processes
	// the escape delimiter. The next character will be stored in the buffer no
	// matter what, and no actions will be dispatched, even if the next
	// character is a token.
	escapeState dfaState = iota

	// valueState is the state encountered after encountering the value
	// delimiter ('='). Being in this state indicates that buffer is no longer
	// storing values for the key.
	valueState dfaState = iota

	// quotedState is the state encountered after the parser encounters a
	// quote. This means that all characters except for the literal escape
	// value will be input into the buffer.
	quotedState dfaState = iota

	// quotedLiteralState is invoked after the parser encounters
	// a `quoteDelimiter` from `quotedState`.
	quotedLiteralState dfaState = iota
)

// PairDelimiter is a type of delimiter that separates different env var key-value pairs
type EnvPairDelimiter = string

const (
	// SemicolonDelimiter (";") is a delimiter between key-value pairs
	SemicolonDelimiter EnvPairDelimiter = ";"

	// CommaDelimiter (",") is a delimiter betewen key-value pairs
	CommaDelimiter EnvPairDelimiter = ","
)

// delimitedString converts environment variables to a particular string
// representation, allowing the user to specify which delimiter to use between
// different environment variable pairs.
func (ev Vars) DelimitedString(separator EnvPairDelimiter) string {
	var serializedPairs []string

	// For each key, value pair, convert it into a "key=value;" pair and
	// continue appending to the output string for each pair
	for k, v := range ev {
		if k != "" {
			pair := quoteDelimiter + escapeString(k, separator) +
				quoteDelimiter + valueDelimiter + quoteDelimiter +
				escapeString(v, separator) + quoteDelimiter
			serializedPairs = append(serializedPairs, pair)
		}
	}
	return strings.Join(serializedPairs, separator)
}

// Parse uses a state machine to parse an input string into the `Vars` type.
// The user can choose which delimiter to use between key-value pairs.
//
// An example of this format:
//
// ENV_VAR_1=VALUE_1;ENV_VAR_2=VALUE_2;
//
// We define the grammar as such (in BNF notation):
// <expr> ::= (<pair> <sep>)* <pair>
// <sep> ::= ';'
//        |  ','
// <pair> ::= <term> = <term>
// <term> ::= <literal>
//         |  "[<literal>|<space>|<escape_quote>]*"
// <literal> ::= [-A-Za-z_0-9]+
// <space> ::= ' '
// <escape_quote> ::= '\"'
func Parse(s string, separator EnvPairDelimiter) (Vars, error) {
	ret := make(Vars)
	var key string
	var value string
	var buffer string
	state := rootState

	// indicates whether the value delimiter has been encountered for the current pair
	valueFlag := false

	for _, c := range s {
		// Having a string is convenient so we can do quality comparisons with
		// the tokens defined in the constants section
		char := string(c)

		// The explanations for each state and what actions should occur in the
		// DFA are found in the comments for each enum
		switch state {
		case rootState:
			// In the case where we have a key-value pair, we want to add that
			// to the map and clear out our buffers
			if char == separator {
				if _, exists := ret[key]; exists {
					return ret, ex.New(fmt.Sprintf("Duplicate keys are not allowed (%s)", key))
				}

				if len(key) == 0 {
					return ret, ex.New("Empty keys are not allowed")
				}

				// This means that we have a term with no '=', which is illegal
				if !valueFlag {
					return ret, ex.New("Expected '='")
				}

				value = buffer
				ret[key] = value

				// clear out the buffers and start over
				buffer = ""
				key = ""
				value = ""
				valueFlag = false
			} else if char == escapeDelimiter {
				state = escapeState
			} else if char == valueDelimiter {
				state = valueState
			} else if char == quoteDelimiter {
				state = quotedState
			} else if unicode.IsSpace(c) {
				continue
			} else {
				buffer += char
			}
		case escapeState:
			buffer += char
			state = rootState
		case valueState:
			if len(buffer) == 0 {
				return ret, ex.New("Empty keys are not allowed")
			}
			key = buffer
			buffer = ""
			valueFlag = true

			if char == quoteDelimiter {
				state = quotedState
			} else {
				if !unicode.IsSpace(c) {
					buffer += char
				}
				state = rootState
			}
		case quotedState:
			if char == escapeDelimiter {
				// ignore the escape and continue
				state = quotedLiteralState
			} else if char == quoteDelimiter {
				state = rootState
			} else {
				buffer += char
			}
		case quotedLiteralState:
			// Escape literal within a quote, goes back to quote mode
			buffer += char
			state = quotedState
		}
	}

	// State 0 is the only valid ending state. If this is not the case, then
	// show the user a parsing error. In the event the input wasn't terminated,
	// we can mitigate by taking the last key-val pair from the buffers.
	switch state {
	case rootState:
		// This handles the case where the key-value pair doesn't have a
		// separator (which is valid grammar). We could go about the option of
		// inserting an extra separator, but that is difficult to do as a
		// preprocessing step because you could have a scenario where there are
		// trailing spaces, or even an escaped ending delimiter.
		if len(buffer) > 0 || len(key) > 0 {
			if !valueFlag {
				return ret, ex.New("Expected '='")
			}
			ret[key] = buffer
		}
	case escapeState:
		return ret, ex.New("Ended input on an escape delimiter ('\\')")
	case valueState:
		return ret, ex.New("Failed to assign a value to some key")
	case quotedState:
		return ret, ex.New("Unclosed quote")
	case quotedLiteralState:
		return ret, ex.New("Ended input on an escape delimiter ('\\')")
	}
	return ret, nil
}

// isToken returns whether a string is a special token that would need to be
// escaped
func isToken(s string, delimiter EnvPairDelimiter) bool {
	switch s {
	case delimiter,
		valueDelimiter,
		quoteDelimiter,
		escapeDelimiter:
		return true
	}
	return false
}

// escapeString takes an string and escapes any special characters so that the
// string can be serialized properly. The user must supply the delimiter used
// to separate key-value pairs.
func escapeString(s string, delimiter EnvPairDelimiter) string {
	var escaped bytes.Buffer
	escaped.Reset()

	for _, r := range s {
		char := string(r)

		if isToken(char, delimiter) {
			escaped.WriteString(escapeDelimiter)
		}
		escaped.WriteRune(r)
	}
	return escaped.String()
}
