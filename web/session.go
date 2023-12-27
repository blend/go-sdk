/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"time"
)

// NewSession returns a new session object.
func NewSession(userID string, sessionID string) *Session {
	return &Session{
		UserID:		userID,
		SessionID:	sessionID,
		CreatedUTC:	time.Now().UTC(),
		State:		map[string]interface{}{},
	}
}

// Session is an active session
type Session struct {
	UserID		string			`json:"userID" yaml:"userID"`
	BaseURL		string			`json:"baseURL" yaml:"baseURL"`
	SessionID	string			`json:"sessionID" yaml:"sessionID"`
	CreatedUTC	time.Time		`json:"createdUTC" yaml:"createdUTC"`
	ExpiresUTC	time.Time		`json:"expiresUTC" yaml:"expiresUTC"`
	UserAgent	string			`json:"userAgent" yaml:"userAgent"`
	RemoteAddr	string			`json:"remoteAddr" yaml:"remoteAddr"`
	State		map[string]interface{}	`json:"state,omitempty" yaml:"state,omitempty"`
}

// WithBaseURL sets the base url.
func (s *Session) WithBaseURL(baseURL string) *Session {
	s.BaseURL = baseURL
	return s
}

// WithUserAgent sets the user agent.
func (s *Session) WithUserAgent(userAgent string) *Session {
	s.UserAgent = userAgent
	return s
}

// WithRemoteAddr sets the remote addr.
func (s *Session) WithRemoteAddr(remoteAddr string) *Session {
	s.RemoteAddr = remoteAddr
	return s
}

// IsExpired returns if the session is expired.
func (s *Session) IsExpired() bool {
	if s.ExpiresUTC.IsZero() {
		return false
	}
	return s.ExpiresUTC.Before(time.Now().UTC())
}

// IsZero returns if the object is set or not.
// It will return true if either the userID or the sessionID are unset.
func (s *Session) IsZero() bool {
	return len(s.UserID) == 0 || len(s.SessionID) == 0
}
