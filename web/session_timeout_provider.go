/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import "time"

// SessionTimeoutProvider returns a new session timeout provider.
func SessionTimeoutProvider(isAbsolute bool, timeout time.Duration) func(*Session) time.Time {
	if isAbsolute {
		return SessionTimeoutProviderAbsolute(timeout)
	}
	return SessionTimeoutProviderRolling(timeout)
}

// SessionTimeoutProviderRolling returns a rolling session timeout.
func SessionTimeoutProviderRolling(timeout time.Duration) func(*Session) time.Time {
	return func(session *Session) time.Time {
		if session.ExpiresUTC.IsZero() {
			return time.Now().UTC().Add(timeout)
		}
		return session.ExpiresUTC.Add(timeout)
	}
}

// SessionTimeoutProviderAbsolute returns an absolute session timeout.
func SessionTimeoutProviderAbsolute(timeout time.Duration) func(*Session) time.Time {
	return func(_ *Session) time.Time {
		return time.Now().UTC().Add(timeout)
	}
}
