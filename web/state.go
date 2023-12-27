/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import "sync"

// State is a provider for a state bag.
type State interface {
	Keys() []string
	Get(key string) interface{}
	Set(key string, value interface{})
	Remove(key string)
	Copy() State
}

// SyncState is the collection of state objects on a context.
type SyncState struct {
	sync.RWMutex
	Values	map[string]interface{}
}

// Keys returns
func (s *SyncState) Keys() (output []string) {
	s.Lock()
	defer s.Unlock()
	if s.Values == nil {
		return
	}

	output = make([]string, len(s.Values))
	var index int
	for key := range s.Values {
		output[index] = key
		index++
	}
	return
}

// Get gets a value.
func (s *SyncState) Get(key string) interface{} {
	s.RLock()
	defer s.RUnlock()
	if s.Values == nil {
		return nil
	}
	return s.Values[key]
}

// Set sets a value.
func (s *SyncState) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	if s.Values == nil {
		s.Values = make(map[string]interface{})
	}
	s.Values[key] = value
}

// Remove removes a key.
func (s *SyncState) Remove(key string) {
	s.Lock()
	defer s.Unlock()
	if s.Values == nil {
		return
	}
	delete(s.Values, key)
}

// Copy creates a new copy of the vars.
func (s *SyncState) Copy() State {
	s.RLock()
	defer s.RUnlock()
	return &SyncState{
		Values: s.Values,
	}
}
