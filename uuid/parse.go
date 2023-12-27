/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package uuid

import "github.com/blend/go-sdk/ex"

// Error Classes
const (
	ErrParseInvalidUUIDInput = ex.Class("parse uuid: existing uuid is invalid")
	// ErrParseEmpty            = ex.Class("parse uuid: input is empty")
	ErrParseInvalidLength    = ex.Class("parse uuid: input is an invalid length")
	ErrParseIllegalCharacter = ex.Class("parse uuid: illegal character")
)

// MustParse parses a uuid and will panic if there is an error.
func MustParse(corpus string) UUID {
	var uuid UUID = make([]byte, 16)
	if err := ParseExisting(&uuid, corpus); err != nil {
		panic(err)
	}
	return uuid
}

// Parse parses a uuidv4 from a given string.
// valid forms are:
// - {xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx}
// - xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// - xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
func Parse(corpus string) (UUID, error) {
	var uuid UUID = make([]byte, 16)
	if err := ParseExisting(&uuid, corpus); err != nil {
		return nil, err
	}
	return uuid, nil
}

// ParseExisting parses into an existing UUID.
func ParseExisting(uuid *UUID, corpus string) error {
	if len(corpus) == 0 {
		return nil // ex.New(ErrParseEmpty)
	}
	if len(corpus)%2 == 1 {
		return ex.New(ErrParseInvalidLength)
	}
	if len(*uuid) != 16 {
		return ex.New(ErrParseInvalidUUIDInput)
	}
	var data = []byte(corpus)
	var c byte
	hex := [2]byte{}
	var hexChar byte
	var isHexChar bool
	var hexIndex, uuidIndex, di int

	for i := 0; i < len(data); i++ {
		c = data[i]
		if c == '{' && i == 0 {
			continue
		}
		if c == '{' {
			return ex.New(ErrParseIllegalCharacter, ex.OptMessagef("at %d: %v", i, string(c)))
		}
		if c == '}' && i != len(data)-1 {
			return ex.New(ErrParseIllegalCharacter, ex.OptMessagef("at %d: %v", i, string(c)))
		}
		if c == '}' {
			continue
		}

		if c == '-' && !(di == 8 || di == 12 || di == 16 || di == 20) {
			return ex.New(ErrParseIllegalCharacter, ex.OptMessagef("at %d: %v", i, string(c)))
		}
		if c == '-' {
			continue
		}

		hexChar, isHexChar = fromHexChar(c)
		if !isHexChar {
			return ex.New(ErrParseIllegalCharacter, ex.OptMessagef("at %d: %v", i, string(c)))
		}

		hex[hexIndex] = hexChar
		if hexIndex == 1 {
			(*uuid)[uuidIndex] = hex[0]<<4 | hex[1]
			uuidIndex++
			hexIndex = 0
		} else {
			hexIndex++
		}
		di++
	}
	if uuidIndex != 16 {
		return ex.New(ErrParseInvalidLength)
	}
	return nil
}

func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}
