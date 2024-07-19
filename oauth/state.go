/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package oauth

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"

	"github.com/blend/go-sdk/ex"
)

// State is the oauth state.
type State struct {
	// Token is a plaintext random token.
	Token string
	// SecureToken is the hashed version of the token.
	// If a key is set, it validates that our app created the oauth state.
	SecureToken string
	// RedirectURI is the redirect uri.
	RedirectURI string
	// Extra includes other state you might need to encode.
	Extra map[string]interface{}
}

// DeserializeState deserializes the oauth state.
func DeserializeState(raw string) (state State, err error) {
	var corpus []byte
	corpus, err = base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		err = ex.New(err)
		return
	}
	buffer := bytes.NewBuffer(corpus)
	if err = gob.NewDecoder(buffer).Decode(&state); err != nil {
		err = ex.New(err)
		return
	}

	return
}

// MustSerializeState serializes a state value but panics if there is an error.
func MustSerializeState(state State) string {
	output, err := SerializeState(state)
	if err != nil {
		panic(err)
	}
	return output
}

// SerializeState serializes the oauth state.
func SerializeState(state State) (output string, err error) {
	buffer := new(bytes.Buffer)
	err = gob.NewEncoder(buffer).Encode(state)
	if err != nil {
		return
	}
	output = base64.RawURLEncoding.EncodeToString(buffer.Bytes())
	return
}
