package configutil

import (
	"context"
	"strconv"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/stringutil"
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
func (p Parser) Bool(ctx context.Context) (*bool, error) {
	value, err := p.Source.String(ctx)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, nil
	}

	parsed, err := stringutil.ParseBool(*value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// Int returns the int value.
func (p Parser) Int(ctx context.Context) (*int, error) {
	value, err := p.Source.String(ctx)
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
func (p Parser) Float64(ctx context.Context) (*float64, error) {
	value, err := p.Source.String(ctx)
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
func (p Parser) Duration(ctx context.Context) (*time.Duration, error) {
	value, err := p.Source.String(ctx)
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
