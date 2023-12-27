/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package ratelimiter

// RateLimiter is a type that can be used as a rate limiter.
type RateLimiter interface {
	// Check returns for a given id `true` if that id is _above_ the rate limit, and false otherwise.
	Check(string) bool
}
