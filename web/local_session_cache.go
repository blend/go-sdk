/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package web

import (
	"context"
	"sync"
)

// NewLocalSessionCache returns a new session cache.
func NewLocalSessionCache() *LocalSessionCache {
	return &LocalSessionCache{
		SessionLock: &sync.Mutex{},
		Sessions:    map[string]*Session{},
	}
}

// LocalSessionCache is a memory cache of sessions.
// It is meant to be used in tests.
type LocalSessionCache struct {
	SessionLock *sync.Mutex
	Sessions    map[string]*Session
}

// Apply applies the local session cache to a given auth manager.
func (lsc *LocalSessionCache) Apply(am *AuthManager) {
	am.FetchHandler = lsc.FetchHandler
	am.PersistHandler = lsc.PersistHandler
	am.RemoveHandler = lsc.RemoveHandler
}

// FetchHandler is a shim to interface with the auth manager.
func (lsc *LocalSessionCache) FetchHandler(_ context.Context, sessionID string) (*Session, error) {
	return lsc.Get(sessionID), nil
}

// PersistHandler is a shim to interface with the auth manager.
func (lsc *LocalSessionCache) PersistHandler(_ context.Context, session *Session) error {
	lsc.Upsert(session)
	return nil
}

// RemoveHandler is a shim to interface with the auth manager.
func (lsc *LocalSessionCache) RemoveHandler(_ context.Context, sessionID string) error {
	lsc.Remove(sessionID)
	return nil
}

// Upsert adds or updates a session to the cache.
func (lsc *LocalSessionCache) Upsert(session *Session) {
	lsc.SessionLock.Lock()
	defer lsc.SessionLock.Unlock()
	lsc.Sessions[session.SessionID] = session
}

// Remove removes a session from the cache.
func (lsc *LocalSessionCache) Remove(sessionID string) {
	lsc.SessionLock.Lock()
	defer lsc.SessionLock.Unlock()
	delete(lsc.Sessions, sessionID)
}

// Get gets a session.
func (lsc *LocalSessionCache) Get(sessionID string) *Session {
	lsc.SessionLock.Lock()
	defer lsc.SessionLock.Unlock()

	if session, hasSession := lsc.Sessions[sessionID]; hasSession {
		return session
	}
	return nil
}
