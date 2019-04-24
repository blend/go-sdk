package assert

import (
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	// RED is the ansi escape code fragment for red.
	RED = "31"
	// BLUE is the ansi escape code fragment for blue.
	BLUE = "94"
	// GREEN is the ansi escape code fragment for green.
	GREEN = "32"
	// YELLOW is the ansi escape code fragment for yellow.
	YELLOW = "33"
	// WHITE is the ansi escape code fragment for white.
	WHITE = "37"
	// GRAY is the ansi escape code fragment for gray.
	GRAY = "90"

	// EMPTY is a constant for the empty (0 length) string.
	EMPTY = ""
)

// Any is a loose type alias to interface{}
type Any = interface{}

// Predicate is a func that returns a bool.
type Predicate func(item Any) bool

//PredicateOfInt is a func that takes an int and returns a bool.
type PredicateOfInt func(item int) bool

// PredicateOfFloat is a func that takes a float64 and returns a bool.
type PredicateOfFloat func(item float64) bool

// PredicateOfString is a func that takes a string and returns a bool.
type PredicateOfString func(item string) bool

// PredicateOfTime is a func that takes a time.Time and returns a bool.
type PredicateOfTime func(item time.Time) bool

// Empty returns an empty assertions class; useful when you want to apply assertions w/o hooking into the testing framework.
func Empty() *Assertions {
	return &Assertions{}
}

// New returns a new instance of `Assertions`.
func New(t *testing.T) *Assertions {
	return &Assertions{
		t:            t,
		timerAbort:   make(chan bool),
		timerAborted: make(chan bool),
	}
}

// Filtered returns a new instance of `Assertions`.
func Filtered(t *testing.T, filter Filter) *Assertions {
	CheckFilter(t, filter)
	return &Assertions{
		filter:       filter,
		t:            t,
		timerAbort:   make(chan bool),
		timerAborted: make(chan bool),
	}
}

// Assertions is the main entry point for using the assertions library.
type Assertions struct {
	output       io.Writer
	filter       Filter
	t            *testing.T
	timerAbort   chan bool
	timerAborted chan bool
}

// WithFilter sets the filter.
func (a *Assertions) WithFilter(filter Filter) *Assertions {
	a.filter = filter
	return a
}

// WithOutput sets the assertions output.
// Error messages will be written to this in addition to the test handler.
func (a *Assertions) WithOutput(w io.Writer) *Assertions {
	a.output = w
	return a
}

// Output returns the underlying output writer.
func (a *Assertions) Output() io.Writer {
	return a.output
}

// assertion represents the actions to take for *each* assertion.
// it is used internally for stats tracking.
func (a *Assertions) assertion() {
	Increment()
}

// NonFatal transitions the assertion into a `NonFatal` assertion; that is, one that will not cause the test to abort if it fails.
// NonFatal assertions are useful when you want to check many properties during a test, but only on an informational basis.
// They will typically return a bool to indicate if the assertion succeeded, or if you should consider the overall
// test to still be a success.
func (a *Assertions) NonFatal() *Optional { //golint you can bite me.
	return &Optional{t: a.t, output: a.output}
}

// NotNil asserts that a reference is not nil.
func (a *Assertions) NotNil(object interface{}, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldNotBeNil(object); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// Nil asserts that a reference is nil.
func (a *Assertions) Nil(object interface{}, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeNil(object); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// Len asserts that a collection has a given length.
func (a *Assertions) Len(collection interface{}, length int, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldHaveLength(collection, length); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// Empty asserts that a collection is empty.
func (a *Assertions) Empty(collection interface{}, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeEmpty(collection); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// NotEmpty asserts that a collection is not empty.
func (a *Assertions) NotEmpty(collection interface{}, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldNotBeEmpty(collection); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// Equal asserts that two objects are deeply equal.
func (a *Assertions) Equal(expected interface{}, actual interface{}, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeEqual(expected, actual); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// ReferenceEqual asserts that two objects are the same reference in memory.
func (a *Assertions) ReferenceEqual(expected interface{}, actual interface{}, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeReferenceEqual(expected, actual); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// NotEqual asserts that two objects are not deeply equal.
func (a *Assertions) NotEqual(expected interface{}, actual interface{}, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldNotBeEqual(expected, actual); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// PanicEqual asserts the panic emitted by an actin equals an expected value.
func (a *Assertions) PanicEqual(expected interface{}, action func(), userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBePanicEqual(expected, action); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// Zero asserts that a value is equal to it's default value.
func (a *Assertions) Zero(value interface{}, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeZero(value); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// NotZero asserts that a value is not equal to it's default value.
func (a *Assertions) NotZero(value interface{}, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeNonZero(value); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// True asserts a boolean is true.
func (a *Assertions) True(object bool, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeTrue(object); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// False asserts a boolean is false.
func (a *Assertions) False(object bool, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeFalse(object); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// InDelta asserts that two floats are within a delta.
//
// The delta is computed by the absolute of the difference betwee `f0` and `f1`
// and testing if that absolute difference is strictly less than `delta`
// if greater, it will fail the assertion, if delta is equal to or greater than difference
// the assertion will pass.
func (a *Assertions) InDelta(f0, f1, delta float64, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeInDelta(f0, f1, delta); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// InTimeDelta asserts that times t1 and t2 are within a delta.
func (a *Assertions) InTimeDelta(t1, t2 time.Time, delta time.Duration, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldBeInTimeDelta(t1, t2, delta); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// FileExists asserts that a file exists at a given filepath on disk.
func (a *Assertions) FileExists(filepath string, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := fileShouldExist(filepath); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// Contains asserts that a substring is present in a corpus.
func (a *Assertions) Contains(corpus, substring string, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldContain(corpus, substring); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// NotContains asserts that a substring is present in a corpus.
func (a *Assertions) NotContains(corpus, substring string, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldNotContain(corpus, substring); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// Any applies a predicate.
func (a *Assertions) Any(target interface{}, predicate Predicate, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldAny(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// AnyOfInt applies a predicate.
func (a *Assertions) AnyOfInt(target []int, predicate PredicateOfInt, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldAnyOfInt(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// AnyOfFloat64 applies a predicate.
func (a *Assertions) AnyOfFloat64(target []float64, predicate PredicateOfFloat, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldAnyOfFloat(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// AnyOfString applies a predicate.
func (a *Assertions) AnyOfString(target []string, predicate PredicateOfString, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldAnyOfString(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// All applies a predicate.
func (a *Assertions) All(target interface{}, predicate Predicate, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldAll(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// AllOfInt applies a predicate.
func (a *Assertions) AllOfInt(target []int, predicate PredicateOfInt, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldAllOfInt(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// AllOfFloat64 applies a predicate.
func (a *Assertions) AllOfFloat64(target []float64, predicate PredicateOfFloat, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldAllOfFloat(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// AllOfString applies a predicate.
func (a *Assertions) AllOfString(target []string, predicate PredicateOfString, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldAllOfString(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// None applies a predicate.
func (a *Assertions) None(target interface{}, predicate Predicate, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldNone(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// NoneOfInt applies a predicate.
func (a *Assertions) NoneOfInt(target []int, predicate PredicateOfInt, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldNoneOfInt(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// NoneOfFloat64 applies a predicate.
func (a *Assertions) NoneOfFloat64(target []float64, predicate PredicateOfFloat, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldNoneOfFloat(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// NoneOfString applies a predicate.
func (a *Assertions) NoneOfString(target []string, predicate PredicateOfString, userMessageComponents ...interface{}) {
	a.assertion()
	if didFail, message := shouldNoneOfString(target, predicate); didFail {
		failNow(a.output, a.t, message, userMessageComponents...)
	}
}

// FailNow forces a test failure (useful for debugging).
func (a *Assertions) FailNow(userMessageComponents ...interface{}) {
	failNow(a.output, a.t, "Fatal Assertion Failed", userMessageComponents...)
}

// StartTimeout starts a timed block.
func (a *Assertions) StartTimeout(timeout time.Duration, userMessageComponents ...interface{}) {
	ticker := time.NewTimer(timeout)
	go func() {
		select {
		case <-ticker.C:
			a.t.Errorf("Timeout Reached")
			a.t.FailNow()
		case <-a.timerAbort:
			a.timerAborted <- true
			return
		}
	}()
}

// EndTimeout marks a timed block as complete.
func (a *Assertions) EndTimeout() {
	a.timerAbort <- true
	<-a.timerAborted
}

// Optional is an assertion type that does not stop a test if an assertion fails, simply outputs the error.
type Optional struct {
	output io.Writer
	t      *testing.T
}

// WithOutput sets an output to capture error output.
func (o *Optional) WithOutput(w io.Writer) *Optional {
	o.output = w
	return o
}

// Output returns the underlying output writer.
func (o *Optional) Output() io.Writer {
	return o.output
}

func (o *Optional) assertion() {
	Increment()
}

// Nil asserts the object is nil.
func (o *Optional) Nil(object interface{}, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeNil(object); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// NotNil asserts the object is not nil.
func (o *Optional) NotNil(object interface{}, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldNotBeNil(object); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// Len asserts that the collection has a specified length.
func (o *Optional) Len(collection interface{}, length int, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldHaveLength(collection, length); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// Empty asserts that a collection is empty.
func (o *Optional) Empty(collection interface{}, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeEmpty(collection); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// NotEmpty asserts that a collection is not empty.
func (o *Optional) NotEmpty(collection interface{}, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldNotBeEmpty(collection); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// Equal asserts that two objects are equal.
func (o *Optional) Equal(expected interface{}, actual interface{}, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeEqual(expected, actual); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// ReferenceEqual asserts that two objects are the same underlying reference in memory.
func (o *Optional) ReferenceEqual(expected interface{}, actual interface{}, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeReferenceEqual(expected, actual); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// NotEqual asserts that two objects are not equal.
func (o *Optional) NotEqual(expected interface{}, actual interface{}, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldNotBeEqual(expected, actual); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// PanicEqual asserts the panic emitted by an actin equals an expected value.
func (o *Optional) PanicEqual(expected interface{}, action func(), userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBePanicEqual(expected, action); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// Zero asserts that a value is the default value.
func (o *Optional) Zero(value interface{}, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeZero(value); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// NotZero asserts that a value is not the default value.
func (o *Optional) NotZero(value interface{}, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeNonZero(value); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// True asserts that a bool is false.
func (o *Optional) True(object bool, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeTrue(object); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// False asserts that a bool is false.
func (o *Optional) False(object bool, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeFalse(object); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// InDelta returns if two float64s are separated by a given delta.
func (o *Optional) InDelta(a, b, delta float64, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeInDelta(a, b, delta); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// InTimeDelta returns if two times are separated by a given delta.
func (o *Optional) InTimeDelta(a, b time.Time, delta time.Duration, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldBeInTimeDelta(a, b, delta); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// FileExists asserts that a file exists on disk at a given filepath.
func (o *Optional) FileExists(filepath string, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := fileShouldExist(filepath); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// Contains checks if a substring is present in a corpus.
func (o *Optional) Contains(corpus, substring string, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldContain(corpus, substring); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// NotContains checks if a substring is not present in a corpus.
func (o *Optional) NotContains(corpus, substring string, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldNotContain(corpus, substring); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// Any applies a predicate.
func (o *Optional) Any(target interface{}, predicate Predicate, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldAny(target, predicate); didFail {
		fail(o.output, o.t, prefixOptional(message), userMessageComponents...)
		return false
	}
	return true
}

// AnyOfInt applies a predicate.
func (o *Optional) AnyOfInt(target []int, predicate PredicateOfInt, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldAnyOfInt(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// AnyOfFloat applies a predicate.
func (o *Optional) AnyOfFloat(target []float64, predicate PredicateOfFloat, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldAnyOfFloat(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// AnyOfString applies a predicate.
func (o *Optional) AnyOfString(target []string, predicate PredicateOfString, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldAnyOfString(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// All applies a predicate.
func (o *Optional) All(target interface{}, predicate Predicate, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldAll(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// AllOfInt applies a predicate.
func (o *Optional) AllOfInt(target []int, predicate PredicateOfInt, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldAllOfInt(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// AllOfFloat applies a predicate.
func (o *Optional) AllOfFloat(target []float64, predicate PredicateOfFloat, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldAllOfFloat(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// AllOfString applies a predicate.
func (o *Optional) AllOfString(target []string, predicate PredicateOfString, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldAllOfString(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// None applies a predicate.
func (o *Optional) None(target interface{}, predicate Predicate, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldNone(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// NoneOfInt applies a predicate.
func (o *Optional) NoneOfInt(target []int, predicate PredicateOfInt, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldNoneOfInt(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// NoneOfFloat applies a predicate.
func (o *Optional) NoneOfFloat(target []float64, predicate PredicateOfFloat, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldNoneOfFloat(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// NoneOfString applies a predicate.
func (o *Optional) NoneOfString(target []string, predicate PredicateOfString, userMessageComponents ...interface{}) bool {
	o.assertion()
	if didFail, message := shouldNoneOfString(target, predicate); didFail {
		fail(o.output, o.t, message, userMessageComponents...)
		return false
	}
	return true
}

// Fail manually injects a failure.
func (o *Optional) Fail(userMessageComponents ...interface{}) {
	fail(o.output, o.t, prefixOptional("Assertion Failed"), userMessageComponents...)
}

// --------------------------------------------------------------------------------
// OUTPUT
// --------------------------------------------------------------------------------

func failNow(w io.Writer, t *testing.T, message string, userMessageComponents ...interface{}) {
	fail(w, t, message, userMessageComponents...)
	if t != nil {
		t.FailNow()
	} else {
		panic(fmt.Errorf(message))
	}
}

func fail(w io.Writer, t *testing.T, message string, userMessageComponents ...interface{}) {
	errorTrace := strings.Join(callerInfo(), "\n\t")

	if len(errorTrace) == 0 {
		errorTrace = "Unknown"
	}

	assertionFailedLabel := color("Assertion Failed!", RED)
	locationLabel := color("Assert Location", GRAY)
	assertionLabel := color("Assertion", GRAY)
	messageLabel := color("Message", GRAY)

	erasure := fmt.Sprintf("\r%s", getClearString())
	userMessage := fmt.Sprint(userMessageComponents...)

	if len(userMessage) > 0 {
		errorFormat := `%s
%s
%s:
	%s
%s:
	%s
%s:
	%s

`
		if t != nil {
			t.Errorf(errorFormat, erasure, assertionFailedLabel, locationLabel, errorTrace, assertionLabel, message, messageLabel, userMessage)
		}
		if w != nil {
			fmt.Fprintf(w, errorFormat, "", assertionFailedLabel, locationLabel, errorTrace, assertionLabel, message, messageLabel, userMessage)
		}

		return

	}
	errorFormat := `%s
%s
%s:
	%s
%s:
	%s

`
	if t != nil {
		t.Errorf(errorFormat, erasure, assertionFailedLabel, locationLabel, errorTrace, assertionLabel, message)
	}
	if w != nil {
		fmt.Fprintf(w, errorFormat, "", assertionFailedLabel, locationLabel, errorTrace, assertionLabel, message)
	}

}

// --------------------------------------------------------------------------------
// ASSERTION LOGIC
// --------------------------------------------------------------------------------

func shouldHaveLength(collection interface{}, length int) (bool, string) {
	if l := getLength(collection); l != length {
		message := shouldBeMultipleMessage(length, l, "Collection should have length")
		return true, message
	}
	return false, EMPTY
}

func shouldNotBeEmpty(collection interface{}) (bool, string) {
	if l := getLength(collection); l == 0 {
		message := "Should not be empty"
		return true, message
	}
	return false, EMPTY
}

func shouldBeEmpty(collection interface{}) (bool, string) {
	if l := getLength(collection); l != 0 {
		message := shouldBeMessage(collection, "Should be empty")
		return true, message
	}
	return false, EMPTY
}

func shouldBeEqual(expected, actual interface{}) (bool, string) {
	if !areEqual(expected, actual) {
		return true, equalMessage(expected, actual)
	}
	return false, EMPTY
}

func shouldBeReferenceEqual(expected, actual interface{}) (bool, string) {
	if !areReferenceEqual(expected, actual) {
		return true, referenceEqualMessage(expected, actual)
	}
	return false, EMPTY
}

func shouldBePanicEqual(expected interface{}, action func()) (bool, string) {
	var actual interface{}
	var didPanic bool
	func() {
		defer func() {
			actual = recover()
			didPanic = actual != nil
		}()
		action()
	}()

	if !didPanic || (didPanic && !areEqual(expected, actual)) {
		return true, panicEqualMessage(didPanic, expected, actual)
	}
	return false, EMPTY
}

func shouldNotBeEqual(expected, actual interface{}) (bool, string) {
	if areEqual(expected, actual) {
		return true, notEqualMessage(expected, actual)
	}
	return false, EMPTY
}

func shouldNotBeNil(object interface{}) (bool, string) {
	if isNil(object) {
		return true, "Should not be nil"
	}
	return false, EMPTY
}

func shouldBeNil(object interface{}) (bool, string) {
	if !isNil(object) {
		return true, shouldBeMessage(object, "Should be nil")
	}
	return false, EMPTY
}

func shouldBeTrue(value bool) (bool, string) {
	if !value {
		return true, "Should be true"
	}
	return false, EMPTY
}

func shouldBeFalse(value bool) (bool, string) {
	if value {
		return true, "Should be false"
	}
	return false, EMPTY
}

func shouldBeZero(value interface{}) (bool, string) {
	if !isZero(value) {
		return true, shouldBeMessage(value, "Should be zero")
	}
	return false, EMPTY
}

func shouldBeNonZero(value interface{}) (bool, string) {
	if isZero(value) {
		return true, "Should be non-zero"
	}
	return false, EMPTY
}

func fileShouldExist(filePath string) (bool, string) {
	_, err := os.Stat(filePath)
	if err != nil {
		pwd, _ := os.Getwd()
		message := fmt.Sprintf("File doesnt exist: %s, `pwd`: %s", filePath, pwd)
		return true, message
	}
	return false, EMPTY
}

func shouldBeInDelta(from, to, delta float64) (bool, string) {
	diff := math.Abs(from - to)
	if diff > delta {
		message := fmt.Sprintf("Absolute difference of %0.5f and %0.5f should be less than %0.5f", from, to, delta)
		return true, message
	}
	return false, EMPTY
}

func shouldBeInTimeDelta(from, to time.Time, delta time.Duration) (bool, string) {
	var diff time.Duration
	if from.After(to) {
		diff = from.Sub(to)
	} else {
		diff = to.Sub(from)
	}
	if diff > delta {
		message := fmt.Sprintf("Delta of %s and %s should be less than %v", from.Format(time.RFC3339), to.Format(time.RFC3339), delta)
		return true, message
	}
	return false, EMPTY
}

func shouldContain(corpus, subString string) (bool, string) {
	if !strings.Contains(corpus, subString) {
		message := fmt.Sprintf("`%s` should contain `%s`", corpus, subString)
		return true, message
	}
	return false, EMPTY
}

func shouldNotContain(corpus, subString string) (bool, string) {
	if strings.Contains(corpus, subString) {
		message := fmt.Sprintf("`%s` should not contain `%s`", corpus, subString)
		return true, message
	}
	return false, EMPTY
}

func shouldAny(target interface{}, predicate Predicate) (bool, string) {
	t := reflect.TypeOf(target)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	v := reflect.ValueOf(target)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if t.Kind() != reflect.Slice {
		return true, "`target` is not a slice"
	}

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface()
		if predicate(obj) {
			return false, EMPTY
		}
	}
	return true, "Predicate did not fire for any element in target"
}

func shouldAnyOfInt(target []int, predicate PredicateOfInt) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(int)
		if predicate(obj) {
			return false, EMPTY
		}
	}
	return true, "Predicate did not fire for any element in target"
}

func shouldAnyOfFloat(target []float64, predicate PredicateOfFloat) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(float64)
		if predicate(obj) {
			return false, EMPTY
		}
	}
	return true, "Predicate did not fire for any element in target"
}

func shouldAnyOfString(target []string, predicate PredicateOfString) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(string)
		if predicate(obj) {
			return false, EMPTY
		}
	}
	return true, "Predicate did not fire for any element in target"
}

func shouldAll(target interface{}, predicate Predicate) (bool, string) {
	t := reflect.TypeOf(target)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	v := reflect.ValueOf(target)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if t.Kind() != reflect.Slice {
		return true, "`target` is not a slice"
	}

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface()
		if !predicate(obj) {
			return true, fmt.Sprintf("Predicate failed for element in target: %#v", obj)
		}
	}
	return false, EMPTY
}

func shouldAllOfInt(target []int, predicate PredicateOfInt) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(int)
		if !predicate(obj) {
			return true, fmt.Sprintf("Predicate failed for element in target: %#v", obj)
		}
	}
	return false, EMPTY
}

func shouldAllOfFloat(target []float64, predicate PredicateOfFloat) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(float64)
		if !predicate(obj) {
			return true, fmt.Sprintf("Predicate failed for element in target: %#v", obj)
		}
	}
	return false, EMPTY
}

func shouldAllOfString(target []string, predicate PredicateOfString) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(string)
		if !predicate(obj) {
			return true, fmt.Sprintf("Predicate failed for element in target: %#v", obj)
		}
	}
	return false, EMPTY
}

func shouldNone(target interface{}, predicate Predicate) (bool, string) {
	t := reflect.TypeOf(target)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	v := reflect.ValueOf(target)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if t.Kind() != reflect.Slice {
		return true, "`target` is not a slice"
	}

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface()
		if predicate(obj) {
			return true, fmt.Sprintf("Predicate passed for element in target: %#v", obj)
		}
	}
	return false, EMPTY
}

func shouldNoneOfInt(target []int, predicate PredicateOfInt) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(int)
		if predicate(obj) {
			return true, fmt.Sprintf("Predicate passed for element in target: %#v", obj)
		}
	}
	return false, EMPTY
}

func shouldNoneOfFloat(target []float64, predicate PredicateOfFloat) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(float64)
		if predicate(obj) {
			return true, fmt.Sprintf("Predicate passed for element in target: %#v", obj)
		}
	}
	return false, EMPTY
}

func shouldNoneOfString(target []string, predicate PredicateOfString) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(string)
		if predicate(obj) {
			return true, fmt.Sprintf("Predicate passed for element in target: %#v", obj)
		}
	}
	return false, EMPTY
}

// --------------------------------------------------------------------------------
// UTILITY
// --------------------------------------------------------------------------------

func prefixOptional(message string) string {
	return "(Non-Fatal) " + message
}

func shouldBeMultipleMessage(expected, actual interface{}, message string) string {
	expectedLabel := color("Expected", WHITE)
	actualLabel := color("Actual", WHITE)

	return fmt.Sprintf(`%s
	%s: 	%#v
	%s: 	%#v`, message, expectedLabel, expected, actualLabel, actual)
}

func shouldBeMessage(object interface{}, message string) string {
	actualLabel := color("Actual", WHITE)
	if err, ok := object.(error); ok {
		return fmt.Sprintf(`%s
	%s: 	%+v`, message, actualLabel, err)
	}
	return fmt.Sprintf(`%s
	%s: 	%#v`, message, actualLabel, object)
}

func notEqualMessage(expected, actual interface{}) string {
	return shouldBeMultipleMessage(expected, actual, "Objects should not be equal")
}

func equalMessage(expected, actual interface{}) string {
	return shouldBeMultipleMessage(expected, actual, "Objects should be equal")
}

func referenceEqualMessage(expected, actual interface{}) string {
	return shouldBeMultipleMessage(expected, actual, "References should be equal")
}

func panicEqualMessage(didPanic bool, expected, actual interface{}) string {
	if !didPanic {
		return "Should have produced a panic"
	}
	return shouldBeMultipleMessage(expected, actual, "Panic from action should equal")
}

func getLength(object interface{}) int {
	if object == nil {
		return 0
	} else if object == "" {
		return 0
	}

	objValue := reflect.ValueOf(object)

	switch objValue.Kind() {
	case reflect.Map:
		fallthrough
	case reflect.Slice, reflect.Chan, reflect.String:
		{
			return objValue.Len()
		}
	}
	return 0
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false
}

func isZero(value interface{}) bool {
	return areEqual(0, value)
}

func areReferenceEqual(expected, actual interface{}) bool {
	if expected == nil && actual == nil {
		return true
	}
	if (expected == nil && actual != nil) || (expected != nil && actual == nil) {
		return false
	}

	return expected == actual
}

func areEqual(expected, actual interface{}) bool {
	if expected == nil && actual == nil {
		return true
	}
	if (expected == nil && actual != nil) || (expected != nil && actual == nil) {
		return false
	}

	actualType := reflect.TypeOf(actual)
	if actualType == nil {
		return false
	}
	expectedValue := reflect.ValueOf(expected)
	if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
		return reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
	}

	return reflect.DeepEqual(expected, actual)
}

func callerInfo() []string {
	pc := uintptr(0)
	file := ""
	line := 0
	ok := false
	name := ""

	callers := []string{}
	for i := 0; ; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			return nil
		}

		if file == "<autogenerated>" {
			break
		}

		parts := strings.Split(file, "/")
		dir := parts[len(parts)-2]
		file = parts[len(parts)-1]
		if dir != "assert" && dir != "mock" && dir != "require" {
			callers = append(callers, fmt.Sprintf("%s:%d", file, line))
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		name = f.Name()

		// Drop the package
		segments := strings.Split(name, ".")
		name = segments[len(segments)-1]
		if isTest(name, "Test") ||
			isTest(name, "Benchmark") ||
			isTest(name, "Example") {
			break
		}
	}

	return callers
}

func color(input string, colorCode string) string {
	return fmt.Sprintf("\033[%s;01m%s\033[0m", colorCode, input)
}

func reflectTypeName(object interface{}) string {
	return reflect.TypeOf(object).Name()
}

func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) { // "Test" is ok
		return true
	}
	rune, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(rune)
}

func getClearString() string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	return strings.Repeat(" ", len(fmt.Sprintf("%s:%d:      ", file, line))+2)
}

func safeExec(action func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	action()
	return
}
