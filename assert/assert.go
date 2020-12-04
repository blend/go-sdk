package assert

import (
	"context"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"
	"unicode"
	"unicode/utf8"
)

// Empty returns an empty assertions handler; useful when you want to apply assertions w/o hooking into the testing framework.
func Empty(opts ...Option) *Assertions {
	a := Assertions{
		OutputFormat: OutputFormatFromEnv(),
		Context:      WithContextID(context.Background(), randomString(8)),
	}
	for _, opt := range opts {
		opt(&a)
	}
	return &a
}

// New returns a new instance of `Assertions`.
func New(t *testing.T, opts ...Option) *Assertions {
	a := Assertions{
		T:            t,
		OutputFormat: OutputFormatFromEnv(),
		Context:      WithContextID(context.Background(), randomString(8)),
	}
	if t != nil {
		a.Context = WithTestName(a.Context, t.Name())
	}
	for _, opt := range opts {
		opt(&a)
	}
	return &a
}

// Assertions is the main entry point for using the assertions library.
type Assertions struct {
	Output       io.Writer
	OutputFormat OutputFormat
	T            *testing.T
	Context      context.Context
	Optional     bool
	Count        int32
}

// Background returns the assertions context.
func (a *Assertions) Background() context.Context {
	return a.Context
}

// assertion represents the actions to take for *each* assertion.
// it is used internally for stats tracking.
func (a *Assertions) assertion() {
	atomic.AddInt32(&a.Count, 1)
}

// NonFatal transitions the assertion into a `NonFatal` assertion; that is, one that will not cause the test to abort if it fails.
// NonFatal assertions are useful when you want to check many properties during a test, but only on an informational basis.
// They will typically return a bool to indicate if the assertion succeeded, or if you should consider the overall
// test to still be a success.
func (a *Assertions) NonFatal() *Assertions { //golint you can bite me.
	return &Assertions{
		T:            a.T,
		Output:       a.Output,
		OutputFormat: a.OutputFormat,
		Optional:     true,
	}
}

func (a *Assertions) fail(message string, userMessageComponents ...interface{}) bool {
	if a.Optional {
		fail(a.Output, a.T, a.OutputFormat, NewFailure(message, userMessageComponents...))
		return false
	}
	failNow(a.Output, a.T, a.OutputFormat, NewFailure(message, userMessageComponents...))
	return false
}

// NotNil asserts that a reference is not nil.
func (a *Assertions) NotNil(object interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNotBeNil(object); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// Nil asserts that a reference is nil.
func (a *Assertions) Nil(object interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeNil(object); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// Len asserts that a collection has a given length.
func (a *Assertions) Len(collection interface{}, length int, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldHaveLength(collection, length); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// Empty asserts that a collection is empty.
func (a *Assertions) Empty(collection interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeEmpty(collection); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// NotEmpty asserts that a collection is not empty.
func (a *Assertions) NotEmpty(collection interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNotBeEmpty(collection); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// Equal asserts that two objects are deeply equal.
func (a *Assertions) Equal(expected interface{}, actual interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeEqual(expected, actual); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// ReferenceEqual asserts that two objects are the same reference in memory.
func (a *Assertions) ReferenceEqual(expected interface{}, actual interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeReferenceEqual(expected, actual); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// NotEqual asserts that two objects are not deeply equal.
func (a *Assertions) NotEqual(expected interface{}, actual interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNotBeEqual(expected, actual); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// PanicEqual asserts the panic emitted by an action equals an expected value.
func (a *Assertions) PanicEqual(expected interface{}, action func(), userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBePanicEqual(expected, action); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// Zero asserts that a value is equal to it's default value.
func (a *Assertions) Zero(value interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeZero(value); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// NotZero asserts that a value is not equal to it's default value.
func (a *Assertions) NotZero(value interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeNonZero(value); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// True asserts a boolean is true.
func (a *Assertions) True(object bool, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeTrue(object); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// False asserts a boolean is false.
func (a *Assertions) False(object bool, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeFalse(object); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// InDelta asserts that two floats are within a delta.
//
// The delta is computed by the absolute of the difference betwee `f0` and `f1`
// and testing if that absolute difference is strictly less than `delta`
// if greater, it will fail the assertion, if delta is equal to or greater than difference
// the assertion will pass.
func (a *Assertions) InDelta(f0, f1, delta float64, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeInDelta(f0, f1, delta); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// InTimeDelta asserts that times t1 and t2 are within a delta.
func (a *Assertions) InTimeDelta(t1, t2 time.Time, delta time.Duration, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldBeInTimeDelta(t1, t2, delta); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// NotInTimeDelta asserts that times t1 and t2 are not within a delta.
func (a *Assertions) NotInTimeDelta(t1, t2 time.Time, delta time.Duration, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNotBeInTimeDelta(t1, t2, delta); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// FileExists asserts that a file exists at a given filepath on disk.
func (a *Assertions) FileExists(filepath string, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldFileExist(filepath); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// Contains asserts that a substring is present in a corpus.
func (a *Assertions) Contains(corpus, substring string, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldContain(corpus, substring); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// NotContains asserts that a substring is present in a corpus.
func (a *Assertions) NotContains(corpus, substring string, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNotContain(corpus, substring); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// Matches returns if a given value matches a given regexp expression.
func (a *Assertions) Matches(expr string, value interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldMatch(expr, value); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// NotMatches returns if a given value does not match a given regexp expression.
func (a *Assertions) NotMatches(expr string, value interface{}, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNotMatch(expr, value); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// Any applies a predicate.
func (a *Assertions) Any(target interface{}, predicate Predicate, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldAny(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// AnyOfInt applies a predicate.
func (a *Assertions) AnyOfInt(target []int, predicate PredicateOfInt, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldAnyOfInt(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// AnyOfFloat64 applies a predicate.
func (a *Assertions) AnyOfFloat64(target []float64, predicate PredicateOfFloat, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldAnyOfFloat(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// AnyOfString applies a predicate.
func (a *Assertions) AnyOfString(target []string, predicate PredicateOfString, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldAnyOfString(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// All applies a predicate.
func (a *Assertions) All(target interface{}, predicate Predicate, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldAll(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// AllOfInt applies a predicate.
func (a *Assertions) AllOfInt(target []int, predicate PredicateOfInt, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldAllOfInt(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// AllOfFloat64 applies a predicate.
func (a *Assertions) AllOfFloat64(target []float64, predicate PredicateOfFloat, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldAllOfFloat(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// AllOfString applies a predicate.
func (a *Assertions) AllOfString(target []string, predicate PredicateOfString, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldAllOfString(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// None applies a predicate.
func (a *Assertions) None(target interface{}, predicate Predicate, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNone(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// NoneOfInt applies a predicate.
func (a *Assertions) NoneOfInt(target []int, predicate PredicateOfInt, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNoneOfInt(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// NoneOfFloat64 applies a predicate.
func (a *Assertions) NoneOfFloat64(target []float64, predicate PredicateOfFloat, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNoneOfFloat(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// NoneOfString applies a predicate.
func (a *Assertions) NoneOfString(target []string, predicate PredicateOfString, userMessageComponents ...interface{}) bool {
	a.assertion()
	if didFail, message := shouldNoneOfString(target, predicate); didFail {
		return a.fail(message, userMessageComponents...)
	}
	return true
}

// FailNow forces a test failure (useful for debugging).
func (a *Assertions) FailNow(userMessageComponents ...interface{}) {
	failNow(a.Output, a.T, a.OutputFormat, NewFailure("Fatal Assertion Failed", userMessageComponents...))
}

// Fail forces a test failure (useful for debugging).
func (a *Assertions) Fail(userMessageComponents ...interface{}) bool {
	fail(a.Output, a.T, a.OutputFormat, NewFailure("Fatal Assertion Failed", userMessageComponents...))
	return true
}

// --------------------------------------------------------------------------------
// OUTPUT
// --------------------------------------------------------------------------------

func failNow(w io.Writer, t *testing.T, outputFormat OutputFormat, failure Failure) {
	fail(w, t, outputFormat, failure)
	if t != nil {
		t.FailNow()
	} else {
		panic(failure)
	}
}

func fail(w io.Writer, t *testing.T, outputFormat OutputFormat, failure Failure) {
	var output string
	switch outputFormat {
	case OutputFormatDefault, OutputFormatText:
		output = fmt.Sprintf("\r%s", getClearString())
		output += failure.Text()
	case OutputFormatJSON:
		output = fmt.Sprintf("\r%s", getLocationString())
		output += failure.JSON()
	default:
		panic(fmt.Errorf("invalid output format: %s", outputFormat))
	}
	if t != nil {
		t.Error(output)
	}
	if w != nil {
		fmt.Fprint(w, output)
	}
}

func callerInfoStrings(frames []stackFrame) []string {
	output := make([]string, len(frames))
	for index := range frames {
		output[index] = frames[index].String()
	}
	return output
}

type stackFrame struct {
	PC       uintptr
	FileFull string
	Dir      string
	File     string
	Name     string
	Line     int
	OK       bool
}

func (sf stackFrame) String() string {
	return fmt.Sprintf("%s:%d", sf.File, sf.Line)
}

func callerInfo() []stackFrame {
	var name string
	var callers []stackFrame
	for i := 0; ; i++ {
		var frame stackFrame
		frame.PC, frame.FileFull, frame.Line, frame.OK = runtime.Caller(i)
		if !frame.OK {
			return nil
		}

		if frame.FileFull == "<autogenerated>" {
			break
		}

		parts := strings.Split(frame.FileFull, "/")
		frame.Dir = parts[len(parts)-2]
		frame.File = parts[len(parts)-1]
		if frame.Dir != "assert" {
			callers = append(callers, frame)
		}

		f := runtime.FuncForPC(frame.PC)
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

func getLocationString() string {
	callers := callerInfo()
	if len(callers) == 0 {
		return ""
	}
	last := callers[len(callers)-1]
	return fmt.Sprintf("%s:%d:      ", last.File, last.Line)
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

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// --------------------------------------------------------------------------------
// ASSERTION LOGIC
// --------------------------------------------------------------------------------

func shouldHaveLength(collection interface{}, length int) (bool, string) {
	if l := getLength(collection); l != length {
		message := shouldBeMultipleMessage(length, l, "Collection should have length")
		return true, message
	}
	return false, ""
}

func shouldNotBeEmpty(collection interface{}) (bool, string) {
	if l := getLength(collection); l == 0 {
		message := "Should not be empty"
		return true, message
	}
	return false, ""
}

func shouldBeEmpty(collection interface{}) (bool, string) {
	if l := getLength(collection); l != 0 {
		message := shouldBeMessage(collection, "Should be empty")
		return true, message
	}
	return false, ""
}

func shouldBeEqual(expected, actual interface{}) (bool, string) {
	if !areEqual(expected, actual) {
		return true, equalMessage(expected, actual)
	}
	return false, ""
}

func shouldBeReferenceEqual(expected, actual interface{}) (bool, string) {
	if !areReferenceEqual(expected, actual) {
		return true, referenceEqualMessage(expected, actual)
	}
	return false, ""
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
	return false, ""
}

func shouldNotBeEqual(expected, actual interface{}) (bool, string) {
	if areEqual(expected, actual) {
		return true, notEqualMessage(expected, actual)
	}
	return false, ""
}

func shouldNotBeNil(object interface{}) (bool, string) {
	if isNil(object) {
		return true, "Should not be nil"
	}
	return false, ""
}

func shouldBeNil(object interface{}) (bool, string) {
	if !isNil(object) {
		return true, shouldBeMessage(object, "Should be nil")
	}
	return false, ""
}

func shouldBeTrue(value bool) (bool, string) {
	if !value {
		return true, "Should be true"
	}
	return false, ""
}

func shouldBeFalse(value bool) (bool, string) {
	if value {
		return true, "Should be false"
	}
	return false, ""
}

func shouldBeZero(value interface{}) (bool, string) {
	if !isZero(value) {
		return true, shouldBeMessage(value, "Should be zero")
	}
	return false, ""
}

func shouldBeNonZero(value interface{}) (bool, string) {
	if isZero(value) {
		return true, "Should be non-zero"
	}
	return false, ""
}

func shouldFileExist(filePath string) (bool, string) {
	_, err := os.Stat(filePath)
	if err != nil {
		pwd, _ := os.Getwd()
		message := fmt.Sprintf("File doesnt exist: %s, `pwd`: %s", filePath, pwd)
		return true, message
	}
	return false, ""
}

func shouldBeInDelta(from, to, delta float64) (bool, string) {
	diff := math.Abs(from - to)
	if diff > delta {
		message := fmt.Sprintf("Absolute difference of %0.5f and %0.5f should be less than %0.5f", from, to, delta)
		return true, message
	}
	return false, ""
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
	return false, ""
}

func shouldNotBeInTimeDelta(from, to time.Time, delta time.Duration) (bool, string) {
	var diff time.Duration
	if from.After(to) {
		diff = from.Sub(to)
	} else {
		diff = to.Sub(from)
	}

	if diff <= delta {
		message := fmt.Sprintf("Delta of %s and %s should be greater than %v", from.Format(time.RFC3339), to.Format(time.RFC3339), delta)
		return true, message
	}
	return false, ""
}

func shouldMatch(pattern string, value interface{}) (bool, string) {
	matched, err := regexp.MatchString(pattern, fmt.Sprint(value))
	if err != nil {
		panic(err)
	}
	if !matched {
		message := fmt.Sprintf("`%v` should match `%s`", value, pattern)
		return true, message
	}
	return false, ""
}

func shouldNotMatch(pattern string, value interface{}) (bool, string) {
	matched, err := regexp.MatchString(pattern, fmt.Sprint(value))
	if err != nil {
		panic(err)
	}
	if matched {
		message := fmt.Sprintf("`%v` should not match `%s`", value, pattern)
		return true, message
	}
	return false, ""
}

func shouldContain(corpus, subString string) (bool, string) {
	if !strings.Contains(corpus, subString) {
		message := fmt.Sprintf("`%s` should contain `%s`", corpus, subString)
		return true, message
	}
	return false, ""
}

func shouldNotContain(corpus, subString string) (bool, string) {
	if strings.Contains(corpus, subString) {
		message := fmt.Sprintf("`%s` should not contain `%s`", corpus, subString)
		return true, message
	}
	return false, ""
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
			return false, ""
		}
	}
	return true, "Predicate did not fire for any element in target"
}

func shouldAnyOfInt(target []int, predicate PredicateOfInt) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(int)
		if predicate(obj) {
			return false, ""
		}
	}
	return true, "Predicate did not fire for any element in target"
}

func shouldAnyOfFloat(target []float64, predicate PredicateOfFloat) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(float64)
		if predicate(obj) {
			return false, ""
		}
	}
	return true, "Predicate did not fire for any element in target"
}

func shouldAnyOfString(target []string, predicate PredicateOfString) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(string)
		if predicate(obj) {
			return false, ""
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
	return false, ""
}

func shouldAllOfInt(target []int, predicate PredicateOfInt) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(int)
		if !predicate(obj) {
			return true, fmt.Sprintf("Predicate failed for element in target: %#v", obj)
		}
	}
	return false, ""
}

func shouldAllOfFloat(target []float64, predicate PredicateOfFloat) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(float64)
		if !predicate(obj) {
			return true, fmt.Sprintf("Predicate failed for element in target: %#v", obj)
		}
	}
	return false, ""
}

func shouldAllOfString(target []string, predicate PredicateOfString) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(string)
		if !predicate(obj) {
			return true, fmt.Sprintf("Predicate failed for element in target: %#v", obj)
		}
	}
	return false, ""
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
	return false, ""
}

func shouldNoneOfInt(target []int, predicate PredicateOfInt) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(int)
		if predicate(obj) {
			return true, fmt.Sprintf("Predicate passed for element in target: %#v", obj)
		}
	}
	return false, ""
}

func shouldNoneOfFloat(target []float64, predicate PredicateOfFloat) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(float64)
		if predicate(obj) {
			return true, fmt.Sprintf("Predicate passed for element in target: %#v", obj)
		}
	}
	return false, ""
}

func shouldNoneOfString(target []string, predicate PredicateOfString) (bool, string) {
	v := reflect.ValueOf(target)

	for x := 0; x < v.Len(); x++ {
		obj := v.Index(x).Interface().(string)
		if predicate(obj) {
			return true, fmt.Sprintf("Predicate passed for element in target: %#v", obj)
		}
	}
	return false, ""
}

// --------------------------------------------------------------------------------
// UTILITY
// --------------------------------------------------------------------------------

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
