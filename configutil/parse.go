package configutil

import (
	"strconv"
	"strings"
	"time"

	"github.com/blend/go-sdk/ex"
)

var (
	_ IntSource      = (*Parser)(nil)
	_ Float64Source  = (*Parser)(nil)
	_ DurationSource = (*Parser)(nil)
)

// Parse returns an int parser.
func Parse(source StringSource) Parser {
	return Parser{Source: source}
}

// Parser parses an int.
type Parser struct {
	Source StringSource
}

// Bool returns the bool value.
func (p Parser) Bool() (*bool, error) {
	value, err := p.Source.String()
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	var parsed bool
	switch strings.ToLower(*value) {
	case "1", "true", "on", "yes":
		parsed = true
	case "0", "false", "off", "no":
		parsed = false
	default:
		return nil, ex.New("invalid bool value", ex.OptMessagef("value: %s", *value))
	}
	return &parsed, nil
}

// Int returns the int value.
func (p Parser) Int() (*int, error) {
	value, err := p.Source.String()
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	parsed, err := strconv.Atoi(*value)
	if err != nil {
		return nil, ex.New(err)
	}
	return &parsed, nil
}

// Float64 returns the float64 value.
func (p Parser) Float64() (*float64, error) {
	value, err := p.Source.String()
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(*value, 64)
	if err != nil {
		return nil, ex.New(err)
	}
	return &parsed, nil
}

// Duration returns a parsed duration value.
func (p Parser) Duration() (*time.Duration, error) {
	value, err := p.Source.String()
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}
	parsed, err := time.ParseDuration(*value)
	if err != nil {
		return nil, ex.New(err)
	}
	return &parsed, nil
}
