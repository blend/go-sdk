/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package selector

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// Parser parses a selector incrementally.
type Parser struct {
	// s stores the string to be tokenized
	s string
	// pos is the position currently tokenized
	pos int
	// m is an optional mark
	m int

	skipValidation bool
}

// Parse does the actual parsing.
func (p *Parser) Parse() (Selector, error) {
	p.s = strings.TrimSpace(p.s)
	if len(p.s) == 0 {
		return Any{}, nil
	}

	var b rune
	var selector, subSelector Selector
	var err error
	var word string
	var op string

	// loop over "clauses"
	// clauses are separated by commas and grouped logically as "ands"
	for {
		// sniff the !haskey form
		b = p.current()

		if b == Bang {
			p.advance() // we aren't going to use the '!'

			// read off the !KEY
			// readWord will leave us on the next non-alpha char
			word, err = p.readWord()
			if err != nil {
				return nil, err
			}

			selector = p.addAnd(selector, p.notHasKey(word)) // add the !KEY term
			if p.done() {
				break
			}

			p.skipToNonWhitespace()
			b = p.current()
			if b != Comma {
				return nil, p.parseError("consecutive not has key terms")
			}

			p.advance()
			continue
		}

		// we're done peeking the first char
		// read the first KEY
		word, err = p.readWord()
		if err != nil {
			return nil, err
		}

		p.mark() // mark to revert if the sniff for the `KEY` form fails

		// sniff if the next character after the word is a comma
		// this indicates it's a "key" form, or existence check on a key
		b = p.skipToNonWhitespace() // the comma is not whitespace
		if b == Comma || p.done() {
			selector = p.addAnd(selector, p.hasKey(word))

			if b == Comma {
				// this is largely a no-op unless we hit a comma
				p.advance()
				// we _have_ to eat the next whitespace
				_ = p.skipToNonWhitespace()
				if p.done() {
					return nil, p.parseError()
				}
			}
			if p.done() {
				break
			}
			continue
		} else {
			p.popMark()
		}

		op, err = p.readOp()
		if err != nil {
			return nil, err
		}

		switch op {
		case OpEquals, OpDoubleEquals:
			subSelector, err = p.equals(word)
		case OpNotEquals:
			subSelector, err = p.notEquals(word)
		case OpIn:
			subSelector, err = p.in(word)
		case OpNotIn:
			subSelector, err = p.notIn(word)
		default:
			return nil, p.parseError("invalid operator")
		}
		if err != nil {
			return nil, err
		}
		selector = p.addAnd(selector, subSelector)

		b = p.skipToNonWhitespace()
		if b == Comma {
			p.advance()
			if p.done() {
				return nil, p.parseError(errExpectedNonEmptyKey)
			}
			p.skipToNonWhitespace()
			continue
		}

		if p.done() {
			break
		}

		// we have a "foo == bar foo" situation
		return nil, p.parseError("keys not separated by comma")
	}

	if !p.skipValidation {
		err = selector.Validate()
		if err != nil {
			return nil, err
		}
	}

	return selector, nil
}

// addAnd starts grouping selectors into a high level `and`, returning the aggregate selector.
func (p *Parser) addAnd(current, next Selector) Selector {
	if current == nil {
		return next
	}
	if typed, isTyped := current.(And); isTyped {
		return append(typed, next)
	}
	return And([]Selector{current, next})
}

func (p *Parser) hasKey(key string) Selector {
	return HasKey(key)
}

func (p *Parser) notHasKey(key string) Selector {
	return NotHasKey(key)
}

func (p *Parser) equals(key string) (Selector, error) {
	value, err := p.readWord()
	if err != nil {
		return nil, err
	}
	return Equals{Key: key, Value: value}, nil
}

func (p *Parser) notEquals(key string) (Selector, error) {
	value, err := p.readWord()
	if err != nil {
		return nil, err
	}
	return NotEquals{Key: key, Value: value}, nil
}

func (p *Parser) in(key string) (Selector, error) {
	csv, err := p.readCSV()
	if err != nil {
		return nil, err
	}
	return In{Key: key, Values: csv}, nil
}

func (p *Parser) notIn(key string) (Selector, error) {
	csv, err := p.readCSV()
	if err != nil {
		return nil, err
	}
	return NotIn{Key: key, Values: csv}, nil
}

// done indicates the cursor is past the usable length of the string.
func (p *Parser) done() bool {
	return p.pos == len(p.s)
}

// mark sets a mark at the current position.
func (p *Parser) mark() {
	p.m = p.pos
}

// popMark moves the cursor back to the previous mark.
func (p *Parser) popMark() {
	if p.m > 0 {
		p.pos = p.m
	}
	p.m = 0
}

// current returns the rune at the current position.
func (p *Parser) current() (r rune) {
	r, _ = utf8.DecodeRuneInString(p.s[p.pos:])
	return
}

// advance moves the cursor forward one rune.
func (p *Parser) advance() {
	if p.pos < len(p.s) {
		_, width := utf8.DecodeRuneInString(p.s[p.pos:])
		p.pos += width
	}
}

// readOp reads a valid operator.
// valid operators include:
// [ =, ==, !=, in, notin ]
// errors if it doesn't read one of the above, or there is another structural issue.
// this will leave the position on the character after the operator
func (p *Parser) readOp() (string, error) {
	// skip preceding whitespace
	p.skipWhiteSpace()

	const (
		stateFirstOpChar = 0
		stateEqual       = 1
		stateBang        = 2
		stateInI         = 3
		stateNotInN      = 4
		stateNotInO      = 5
		stateNotInT      = 6
		stateNotInI      = 7
	)

	var state int
	var ch rune
	var op []rune
	for {
		if p.done() {
			return "", p.parseError("invalid operator")
		}

		ch = p.current()

		switch state {
		case stateFirstOpChar: // initial state, determine what op we're reading for
			if ch == Equal {
				state = stateEqual
				break
			}
			if ch == Bang {
				state = stateBang
				break
			}
			if ch == 'i' {
				state = stateInI
				break
			}
			if ch == 'n' {
				state = stateNotInN
				break
			}

			return "", p.parseError("invalid operator")

		case stateEqual:
			if p.isWhitespace(ch) || isAlpha(ch) || ch == Comma {
				return string(op), nil
			}
			if ch == Equal {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}

			return "", p.parseError("invalid operator")

		case stateBang:
			if ch == Equal {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}

			return "", p.parseError("invalid operator")

		case stateInI:
			if ch == 'n' {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}

			return "", p.parseError("invalid operator")

		case stateNotInN:
			if ch == 'o' {
				state = stateNotInO
				break
			}

			return "", p.parseError("invalid operator")

		case stateNotInO:
			if ch == 't' {
				state = stateNotInT
				break
			}

			return "", p.parseError("invalid operator")

		case stateNotInT:
			if ch == 'i' {
				state = stateNotInI
				break
			}

			return "", p.parseError("invalid operator")

		case stateNotInI:
			if ch == 'n' {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}

			return "", p.parseError("invalid operator")
		}

		op = append(op, ch)
		p.advance()
	}
}

// readWord skips whitespace, then reads a word until whitespace or a token.
// it will leave the cursor on the next char after the word, i.e. the space or token.
func (p *Parser) readWord() (string, error) {
	p.skipWhiteSpace()

	var word []rune
	var ch rune
	for {
		if p.done() {
			break
		}

		ch = p.current()
		if isWhitespace(ch) ||
			ch == Comma ||
			isOperatorSymbol(ch) {
			break
		}

		word = append(word, ch)
		p.advance()
	}

	if len(word) == 0 {
		return "", p.parseError(errExpectedNonEmptyKey)
	}

	return string(word), nil
}

// readCSV reads an array of strings in csv form.
// it expects to start just before the first `(` and
// will read until just past the closing `)`
func (p *Parser) readCSV() (results []string, err error) {
	// skip preceding whitespace
	p.skipWhiteSpace()

	const (
		stateBeforeParens              = 0
		stateWord                      = 1
		stateWhitespaceAfterOpenParens = 2
		stateWhitespaceAfterComma      = 3
		stateWhitespaceAfterWord       = 4
	)

	var word []rune
	var ch rune
	var state int

	for {
		if p.done() {
			results = nil
			err = p.parseError("csv; expects close parenthesis")
			// err = ErrInvalidSelector
			return
		}

		ch = p.current()

		switch state {
		case stateBeforeParens:
			if ch == OpenParens {
				state = stateWhitespaceAfterOpenParens
				p.advance()
				continue
			}

			// not open parens, bail
			err = p.parseError("csv; expects open parenthesis")
			results = nil
			return

		case stateWord:

			if ch == Comma {
				if len(word) > 0 {
					results = append(results, string(word))
					word = nil
				}

				// the symbol is the comma
				state = stateWhitespaceAfterComma
				p.advance()
				continue
			}

			if ch == CloseParens {
				if len(word) > 0 {
					results = append(results, string(word))
				}
				p.advance()
				return
			}

			if p.isWhitespace(ch) {
				if len(word) > 0 {
					results = append(results, string(word))
					word = nil
				}

				state = stateWhitespaceAfterWord
				p.advance()
				continue
			}

			if !p.isValidValue(ch) {
				err = p.parseError("csv; word contains invalid characters")
				results = nil
				return
			}

			word = append(word, ch)
			p.advance()
			continue

		case stateWhitespaceAfterOpenParens, stateWhitespaceAfterComma:
			if p.isWhitespace(ch) {
				p.advance()
				continue
			}
			if isAlpha(ch) {
				state = stateWord
				continue
			}
			if ch == Comma {
				p.advance()
				state = stateWhitespaceAfterComma
				continue
			}
			if ch == CloseParens {
				p.advance()
				return // exit reading the csv
			}

			if state == stateWhitespaceAfterOpenParens {
				err = p.parseError("csv; invalid characters after '('")
				return
			}
			err = p.parseError("csv; invalid characters after ','")
			return

		case stateWhitespaceAfterWord:

			if ch == CloseParens {
				if len(word) > 0 {
					results = append(results, string(word))
				}
				p.advance()
				return
			}

			if p.isWhitespace(ch) {
				p.advance()
				continue
			}

			if ch == Comma {
				state = stateWhitespaceAfterComma
				p.advance()
				continue
			}

			err = p.parseError("csv; consecutive whitespace separated words without a comma")
			results = nil
			return
		}
	}
}

func (p *Parser) skipWhiteSpace() {
	var ch rune
	for {
		if p.done() {
			return
		}
		ch = p.current()
		if !p.isWhitespace(ch) {
			return
		}
		p.advance()
	}
}

func (p *Parser) skipToNonWhitespace() (ch rune) {
	for {
		if p.done() {
			return
		}
		ch = p.current()
		if ch == Comma || !p.isWhitespace(ch) {
			return
		}
		p.advance()
	}
}

// isWhitespace returns true if the rune is a space, tab, or newline.
func (p *Parser) isWhitespace(ch rune) bool {
	return ch == Space || ch == Tab || ch == CarriageReturn || ch == NewLine
}

func (p *Parser) isValidValue(ch rune) bool {
	return isAlpha(ch) || isNameSymbol(ch)
}

func (p *Parser) parseError(message ...interface{}) error {
	return &ParseError{
		Err:      ErrInvalidSelector,
		Input:    p.s,
		Position: p.pos,
		Message:  fmt.Sprint(message...),
	}
}
