package selector

import (
	"strings"
	"unicode/utf8"
)

const (
	// OpEquals is an operator.
	OpEquals = "="
	// OpDoubleEquals is an operator.
	OpDoubleEquals = "=="
	// OpNotEquals is an operator.
	OpNotEquals = "!="
	// OpIn is an operator.
	OpIn = "in"
	// OpNotIn is an operator.
	OpNotIn = "notin"
)

// Parser parses a selector incrementally.
type Parser struct {
	// s stores the string to be tokenized
	s string
	// pos is the position currently tokenized
	pos int
	// m is an optional mark
	m int

	skipValidation  bool
	permittedValues map[rune]bool
	selectorOptions []SelectorOption
}

// Parse does the actual parsing.
func (p *Parser) Parse() (Selector, error) {
	p.s = strings.TrimSpace(p.s)
	if len(p.s) == 0 {
		return Any{}, nil
	}

	var b rune
	var selector Selector
	var err error
	var op string

	// loop over "clauses"
	// clauses are separated by commas and grouped logically as "ands"
	for {

		// sniff the !haskey form
		b = p.current()
		if b == Bang {
			p.advance() // we aren't going to use the '!'
			selector = p.addAnd(selector, p.notHasKey(p.readWord()))
			if p.done() {
				break
			}
			continue
		}

		// we're done peeking the first char
		key := p.readWord()

		p.mark()

		// check if the next character after the word is a comma
		// this indicates it's a "key" form, or existence check on a key
		b = p.skipToComma()
		if b == Comma || p.isTerminator(b) || p.done() {
			selector = p.addAnd(selector, p.hasKey(key))
			p.advance()
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

		var subSelector Selector
		switch op {
		case OpEquals, OpDoubleEquals:
			subSelector, err = p.equals(key)
			if err != nil {
				return nil, err
			}
			selector = p.addAnd(selector, subSelector)
		case OpNotEquals:
			subSelector, err = p.notEquals(key)
			if err != nil {
				return nil, err
			}
			selector = p.addAnd(selector, subSelector)
		case OpIn:
			subSelector, err = p.in(key)
			if err != nil {
				return nil, err
			}
			selector = p.addAnd(selector, subSelector)
		case OpNotIn:
			subSelector, err = p.notIn(key)
			if err != nil {
				return nil, err
			}
			selector = p.addAnd(selector, subSelector)
		default:
			return nil, ErrInvalidOperator
		}

		b = p.skipToComma()
		if b == Comma {
			p.advance()
			if p.done() {
				break
			}
			continue
		}

		// these two are effectively the same
		if p.isTerminator(b) || p.done() {
			break
		}

		return nil, ErrInvalidSelector
	}

	if !p.skipValidation {
		err = selector.Validate(p.selectorOptions...)
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
	if typed, isTyped := current.(*And); isTyped {
		typed.selectors = append(typed.selectors, next)
		return typed
	}
	return &And{selectors: []Selector{current, next}}
}

func (p *Parser) hasKey(key string) Selector {
	return HasKey(key)
}

func (p *Parser) notHasKey(key string) Selector {
	return NotHasKey(key)
}

func (p *Parser) equals(key string) (Selector, error) {
	value := p.readWord()
	return &Equals{Key: key, Value: value}, nil
}

func (p *Parser) notEquals(key string) (Selector, error) {
	value := p.readWord()
	return &NotEquals{Key: key, Value: value}, nil
}

func (p *Parser) in(key string) (Selector, error) {
	csv, err := p.readCSV()
	if err != nil {
		return nil, err
	}
	return &In{Key: key, Values: csv}, nil
}

func (p *Parser) notIn(key string) (Selector, error) {
	csv, err := p.readCSV()
	if err != nil {
		return nil, err
	}
	return &NotIn{Key: key, Values: csv}, nil
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

// read returns the rune currently lexed, and advances the position.
func (p *Parser) read() (r rune) {
	var width int
	if p.pos < len(p.s) {
		r, width = utf8.DecodeRuneInString(p.s[p.pos:])
		p.pos += width
	}
	return r
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

// prev moves the cursor back a rune.
func (p *Parser) prev() {
	if p.pos > 0 {
		p.pos--
	}
}

// readOp reads a valid operator.
// valid operators include:
// [ =, ==, !=, in, notin ]
// errors if it doesn't read one of the above, or there is another structural issue.
func (p *Parser) readOp() (string, error) {
	// skip preceding whitespace
	p.skipWhiteSpace()

	var state int
	var ch rune
	var op []rune
	for {
		ch = p.current()

		switch state {
		case 0: // initial state, determine what op we're reading for
			if ch == Equal {
				state = 1
				break
			}
			if ch == Bang {
				state = 2
				break
			}
			if ch == 'i' {
				state = 6
				break
			}
			if ch == 'n' {
				state = 7
				break
			}
			return "", ErrInvalidOperator
		case 1: // =
			if p.isWhitespace(ch) || p.isAlpha(ch) || ch == Comma || isPermitted(ch, p.permittedValues) {
				return string(op), nil
			}
			if ch == Equal {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 2: // !
			if ch == Equal {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 6: // in
			if ch == 'n' {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		case 7: // o
			if ch == 'o' {
				state = 8
				break
			}
			return "", ErrInvalidOperator
		case 8: // t
			if ch == 't' {
				state = 9
				break
			}
			return "", ErrInvalidOperator
		case 9: // i
			if ch == 'i' {
				state = 10
				break
			}
			return "", ErrInvalidOperator
		case 10: // n
			if ch == 'n' {
				op = append(op, ch)
				p.advance()
				return string(op), nil
			}
			return "", ErrInvalidOperator
		}

		op = append(op, ch)
		p.advance()

		if p.done() {
			return string(op), nil
		}
	}
}

// readWord skips whitespace, then reads a word until whitespace or a token.
// it will leave the cursor on the next char after the word, i.e. the space or token.
func (p *Parser) readWord() string {
	// skip preceding whitespace
	p.skipWhiteSpace()

	var word []rune
	var ch rune
	for {
		ch = p.current()

		if p.isWhitespace(ch) {
			return string(word)
		}
		if p.isSpecialSymbol(ch) {
			return string(word)
		}

		word = append(word, ch)
		p.advance()

		if p.done() {
			return string(word)
		}
	}
}

func (p *Parser) readCSV() (results []string, err error) {
	// skip preceding whitespace
	p.skipWhiteSpace()

	var word []rune
	var ch rune
	var state int

	for {
		ch = p.current()

		if p.done() {
			err = ErrInvalidSelector
			return
		}

		switch state {
		case 0: // leading paren
			if ch == OpenParens {
				state = 2 // spaces or alphas
				p.advance()
				continue
			}
			// not open parens, bail
			err = ErrInvalidSelector
			return
		case 1: // alphas (in word)

			if ch == Comma {
				if len(word) > 0 {
					results = append(results, string(word))
					word = nil
				}
				state = 2 // from comma
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
				state = 3
				p.advance()
				continue
			}

			if !p.isValidValue(ch) {
				err = ErrInvalidSelector
				return
			}

			word = append(word, ch)
			p.advance()
			continue

		case 2: //whitespace after symbol

			if ch == CloseParens {
				p.advance()
				return
			}

			if p.isWhitespace(ch) {
				p.advance()
				continue
			}

			if ch == Comma {
				p.advance()
				continue
			}

			if p.isAlpha(ch) {
				state = 1
				continue
			}

			err = ErrInvalidSelector
			return

		case 3: //whitespace after alpha

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
				if len(word) > 0 {
					results = append(results, string(word))
					word = nil
				}
				p.advance()
				state = 2
				continue
			}

			err = ErrInvalidSelector
			return

		}
	}
}

func (p *Parser) skipWhiteSpace() {
	if p.done() {
		return
	}
	var ch rune
	for {
		ch = p.current()
		if !p.isWhitespace(ch) {
			return
		}
		p.advance()
		if p.done() {
			return
		}
	}
}

func (p *Parser) skipToComma() (ch rune) {
	if p.done() {
		return
	}
	for {
		ch = p.current()
		if ch == Comma {
			return
		}
		if !p.isWhitespace(ch) {
			return
		}
		p.advance()
		if p.done() {
			return
		}
	}
}

// isWhitespace returns true if the rune is a space, tab, or newline.
func (p *Parser) isWhitespace(ch rune) bool {
	return ch == Space || ch == Tab || ch == CarriageReturn || ch == NewLine
}

// isSpecialSymbol returns if the ch is on the selector symbol list.
func (p *Parser) isSpecialSymbol(ch rune) bool {
	return isSelectorSymbol(ch)
}

// isTerminator returns if we've reached the end of the string
func (p *Parser) isTerminator(ch rune) bool {
	return ch == 0
}

func (p *Parser) isAlpha(ch rune) bool {
	return isAlpha(ch)
}

func (p *Parser) isValidValue(ch rune) bool {
	return isAlpha(ch) || isNameSymbol(ch)
}
