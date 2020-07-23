package assert

import "time"

// Predicate is a func that returns a bool.
type Predicate func(item interface{}) bool

//PredicateOfInt is a func that takes an int and returns a bool.
type PredicateOfInt func(item int) bool

// PredicateOfFloat is a func that takes a float64 and returns a bool.
type PredicateOfFloat func(item float64) bool

// PredicateOfString is a func that takes a string and returns a bool.
type PredicateOfString func(item string) bool

// PredicateOfTime is a func that takes a time.Time and returns a bool.
type PredicateOfTime func(item time.Time) bool
