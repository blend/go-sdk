/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package oauth

// StateOption is an option for state objects
type StateOption func(*State)

// OptStateSecureToken sets the secure token on the state.
func OptStateSecureToken(secureToken string) StateOption {
	return func(s *State) {
		s.SecureToken = secureToken
	}
}

// OptStateRedirectURI sets the redirect uri on the stae.
func OptStateRedirectURI(redirectURI string) StateOption {
	return func(s *State) {
		s.RedirectURI = redirectURI
	}
}

// OptStateExtra sets the redirect uri on the stae.
func OptStateExtra(key string, value interface{}) StateOption {
	return func(s *State) {
		if s.Extra == nil {
			s.Extra = make(map[string]interface{})
		}
		s.Extra[key] = value
	}
}
