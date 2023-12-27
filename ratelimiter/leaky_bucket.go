/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ratelimiter

import (
	"time"
)

var (
	_ RateLimiter = (*LeakyBucket)(nil)
)

// NewLeakyBucket returns a new token bucket rate limiter.
// The rate is formed by `numActions` and `quantum`; the resulting rate is numActions/quantum.
func NewLeakyBucket(numActions int, quantum time.Duration) *LeakyBucket {
	return &LeakyBucket{
		NumActions: numActions,
		Quantum:    quantum,
		Tokens:     make(map[string]*Token),
		Now:        func() time.Time { return time.Now().UTC() },
	}
}

// LeakyBucket implements the token bucket rate limiting algorithm.
type LeakyBucket struct {
	NumActions int
	Quantum    time.Duration
	Tokens     map[string]*Token
	Now        func() time.Time
}

// Check returns true if an id has exceeded the rate limit, and false otherwise.
func (lb *LeakyBucket) Check(id string) bool {
	now := lb.Now()

	if lb.Tokens == nil {
		lb.Tokens = make(map[string]*Token)
	}

	token, ok := lb.Tokens[id]
	if !ok {
		lb.Tokens[id] = &Token{Count: 1, Last: now}
		return false
	}

	elapsed := now.Sub(token.Last) // how long since the last call
	// uint64 is used here because of how mantissa bones these calculations in float64
	leakBy := uint64(lb.NumActions) * (uint64(elapsed) / uint64(lb.Quantum))

	token.Count = token.Count - float64(leakBy) // remove by the rate per quantum
	if token.Count < 0 {
		token.Count = 0
	}
	token.Last = now
	token.Count++

	return token.Count >= float64(lb.NumActions)
}

// Token is an individual id's work.
type Token struct {
	Count float64   // the rate adjusted count; initialize at max*rate, remove rate tokens per call
	Last  time.Time // last is used to calculate the elapsed, and subsequently the rate
}
