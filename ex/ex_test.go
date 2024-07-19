/*

Copyright (c) 2024 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package ex

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/blend/go-sdk/assert"
)

func TestNewOfString(t *testing.T) {
	a := assert.New(t)
	ex := As(New("this is a test"))
	a.Equal("this is a test", fmt.Sprintf("%v", ex))
	a.NotNil(ex.StackTrace)
	a.Nil(ex.Inner)
}

func TestNewOfError(t *testing.T) {
	a := assert.New(t)

	err := errors.New("This is an error")
	wrappedErr := New(err)
	a.NotNil(wrappedErr)
	typedWrapped := As(wrappedErr)
	a.NotNil(typedWrapped)
	a.Equal("This is an error", fmt.Sprintf("%v", typedWrapped))
}

func TestNewOfException(t *testing.T) {
	a := assert.New(t)
	ex := New(Class("This is an exception"))
	wrappedEx := New(ex)
	a.NotNil(wrappedEx)
	typedWrappedEx := As(wrappedEx)
	a.Equal("This is an exception", fmt.Sprintf("%v", typedWrappedEx))
	a.Equal(ex, typedWrappedEx)
}

func TestNewOfNil(t *testing.T) {
	a := assert.New(t)

	shouldBeNil := New(nil)
	a.Nil(shouldBeNil)
	a.Equal(nil, shouldBeNil)
	a.True(nil == shouldBeNil)
}

func TestNewOfTypedNil(t *testing.T) {
	a := assert.New(t)

	var nilError error
	a.Nil(nilError)
	a.Equal(nil, nilError)

	shouldBeNil := New(nilError)
	a.Nil(shouldBeNil)
	a.True(shouldBeNil == nil)
}

func TestNewOfReturnedNil(t *testing.T) {
	a := assert.New(t)

	returnsNil := func() error {
		return nil
	}

	shouldBeNil := New(returnsNil())
	a.Nil(shouldBeNil)
	a.True(shouldBeNil == nil)

	returnsTypedNil := func() error {
		return New(nil)
	}

	shouldAlsoBeNil := returnsTypedNil()
	a.Nil(shouldAlsoBeNil)
	a.True(shouldAlsoBeNil == nil)
}

func TestError(t *testing.T) {
	a := assert.New(t)

	ex := New(Class("this is a test"))
	message := ex.Error()
	a.NotEmpty(message)
}

func TestErrorOptions(t *testing.T) {
	a := assert.New(t)

	ex := New(Class("this is a test"), OptMessage("foo"))
	message := ex.Error()
	a.NotEmpty(message)

	typed := As(ex)
	a.NotNil(typed)
	a.Equal("foo", typed.Message)
}

func TestCallers(t *testing.T) {
	a := assert.New(t)

	callStack := func() StackTrace { return Callers(DefaultStartDepth) }()

	a.NotNil(callStack)
	callstackStr := callStack.String()
	a.True(strings.Contains(callstackStr, "TestCallers"), callstackStr)
}

func TestExceptionFormatters(t *testing.T) {
	assert := assert.New(t)

	// test the "%v" formatter with just the exception class.
	class := &Ex{Class: Class("this is a test")}
	assert.Equal("this is a test", fmt.Sprintf("%v", class))

	classAndMessage := &Ex{Class: Class("foo"), Message: "bar"}
	assert.Equal("foo; bar", fmt.Sprintf("%v", classAndMessage))
}

func TestSerializeExamples(t *testing.T) {
	it := assert.New(t)

	n := (*Ex)(nil)
	e := Ex{}
	c1 := Ex{Class: Class("cls")}
	c2 := Ex{Class: fmt.Errorf("cls")}
	m := Ex{Message: "msg"}
	i1 := Ex{Inner: Class("inner")}
	i2 := Ex{Inner: &Ex{Class: Class("inner")}}
	i3 := Ex{Inner: &Ex{Message: "inner"}}
	i4 := Ex{Inner: fmt.Errorf("inner")}
	s := Ex{StackTrace: StackStrings{"stack"}}
	mc1 := Ex{Message: m.Message, Class: c1.Class}
	mc2 := Ex{Message: m.Message, Class: c2.Class}
	msc1 := Ex{Message: m.Message, Class: c1.Class, StackTrace: s.StackTrace}
	msc2 := Ex{Message: m.Message, Class: c2.Class, StackTrace: s.StackTrace}

	// Error ((*Ex).Error)
	// it.Equal("<nil>", n.Error()) // panics
	// it.Equal("", e.Error()) // panics
	it.Equal("cls", c1.Error())
	it.Equal("cls", c2.Error())
	// it.Equal(" msg", m.Error()) // panics
	// it.Equal("", i1.Error()) // panics
	// it.Equal("", i2.Error()) // panics
	// it.Equal("", i3.Error()) // panics
	// it.Equal("", i4.Error()) // panics
	// it.Equal(" \nstack", s.Error()) // panics
	it.Equal("cls", mc1.Error())
	it.Equal("cls", mc2.Error())
	it.Equal("cls", msc1.Error())
	it.Equal("cls", msc2.Error())

	// Stringer ((*Ex).String)
	// it.Equal("<nil>", n.String()) // panics
	it.Equal("", e.String())
	it.Equal("cls", c1.String())
	it.Equal("cls", c2.String())
	it.Equal(" msg", m.String())
	it.Equal("", i1.String())
	it.Equal("", i2.String())
	it.Equal("", i3.String())
	it.Equal("", i4.String())
	it.Equal(" \nstack", s.String())
	it.Equal("cls msg", mc1.String())
	it.Equal("cls msg", mc2.String())
	it.Equal("cls msg \nstack", msc1.String())
	it.Equal("cls msg \nstack", msc2.String())

	// JSON ((*Ex).MarshalJSON)
	var b []byte
	var err error
	// b, err = n.MarshalJSON() // panics
	// b, err = e.MarshalJSON() // panics
	b, err = c1.MarshalJSON()
	it.Nil(err)
	it.Equal(`{"Class":"cls","Message":""}`, string(b))
	b, err = c2.MarshalJSON()
	it.Nil(err)
	it.Equal(`{"Class":"cls","Message":""}`, string(b))
	// b, err = m.MarshalJSON() // panics
	// b, err = i1.MarshalJSON() // panics
	// b, err = i2.MarshalJSON() // panics
	// b, err = i3.MarshalJSON() // panics
	// b, err = i4.MarshalJSON() // panics
	// b, err = s.MarshalJSON() // panics
	b, err = mc1.MarshalJSON()
	it.Nil(err)
	it.Equal(`{"Class":"cls","Message":"msg"}`, string(b))
	b, err = mc2.MarshalJSON()
	it.Nil(err)
	it.Equal(`{"Class":"cls","Message":"msg"}`, string(b))
	b, err = msc1.MarshalJSON()
	it.Nil(err)
	it.Equal(`{"Class":"cls","Message":"msg","StackTrace":["stack"]}`, string(b))
	b, err = msc2.MarshalJSON()
	it.Nil(err)
	it.Equal(`{"Class":"cls","Message":"msg","StackTrace":["stack"]}`, string(b))

	// Format: Ex -> v
	it.Equal("<nil>", fmt.Sprintf("%v", n))
	it.Equal("{<nil>  <nil> <nil>}", fmt.Sprintf("%v", e))
	it.Equal("{cls  <nil> <nil>}", fmt.Sprintf("%v", c1))
	it.Equal("{cls  <nil> <nil>}", fmt.Sprintf("%v", c2))
	it.Equal("{<nil> msg <nil> <nil>}", fmt.Sprintf("%v", m))
	it.Equal("{<nil>  inner <nil>}", fmt.Sprintf("%v", i1))
	it.Equal("{<nil>  inner <nil>}", fmt.Sprintf("%v", i2))
	it.Equal("{<nil>  ; inner <nil>}", fmt.Sprintf("%v", i3))
	it.Equal("{<nil>  inner <nil>}", fmt.Sprintf("%v", i4))
	it.Equal("{<nil>  <nil> \nstack}", fmt.Sprintf("%v", s))
	it.Equal("{cls msg <nil> <nil>}", fmt.Sprintf("%v", mc1))
	it.Equal("{cls msg <nil> <nil>}", fmt.Sprintf("%v", mc2))
	it.Equal("{cls msg <nil> \nstack}", fmt.Sprintf("%v", msc1))
	it.Equal("{cls msg <nil> \nstack}", fmt.Sprintf("%v", msc2))

	// Format: *Ex -> v
	it.Equal("", fmt.Sprintf("%v", &e))
	it.Equal("cls", fmt.Sprintf("%v", &c1))
	it.Equal("cls", fmt.Sprintf("%v", &c2))
	it.Equal("; msg", fmt.Sprintf("%v", &m))
	it.Equal("\ninner", fmt.Sprintf("%v", &i1))
	it.Equal("\ninner", fmt.Sprintf("%v", &i2))
	it.Equal("\n; inner", fmt.Sprintf("%v", &i3))
	it.Equal("\ninner", fmt.Sprintf("%v", &i4))
	it.Equal("", fmt.Sprintf("%v", &s))
	it.Equal("cls; msg", fmt.Sprintf("%v", &mc1))
	it.Equal("cls; msg", fmt.Sprintf("%v", &mc2))
	it.Equal("cls; msg", fmt.Sprintf("%v", &msc1))
	it.Equal("cls; msg", fmt.Sprintf("%v", &msc2))

	// Format: Ex -> +v
	it.Equal("<nil>", fmt.Sprintf("%+v", n))
	it.Equal("{Class:<nil> Message: Inner:<nil> StackTrace:<nil>}", fmt.Sprintf("%+v", e))
	it.Equal("{Class:cls Message: Inner:<nil> StackTrace:<nil>}", fmt.Sprintf("%+v", c1))
	it.Equal("{Class:cls Message: Inner:<nil> StackTrace:<nil>}", fmt.Sprintf("%+v", c2))
	it.Equal("{Class:<nil> Message:msg Inner:<nil> StackTrace:<nil>}", fmt.Sprintf("%+v", m))
	it.Equal("{Class:<nil> Message: Inner:inner StackTrace:<nil>}", fmt.Sprintf("%+v", i1))
	it.Equal("{Class:<nil> Message: Inner:inner StackTrace:<nil>}", fmt.Sprintf("%+v", i2))
	it.Equal("{Class:<nil> Message: Inner:; inner StackTrace:<nil>}", fmt.Sprintf("%+v", i3))
	it.Equal("{Class:<nil> Message: Inner:inner StackTrace:<nil>}", fmt.Sprintf("%+v", i4))
	it.Equal("{Class:<nil> Message: Inner:<nil> StackTrace:\nstack}", fmt.Sprintf("%+v", s))
	it.Equal("{Class:cls Message:msg Inner:<nil> StackTrace:<nil>}", fmt.Sprintf("%+v", mc1))
	it.Equal("{Class:cls Message:msg Inner:<nil> StackTrace:<nil>}", fmt.Sprintf("%+v", mc2))
	it.Equal("{Class:cls Message:msg Inner:<nil> StackTrace:\nstack}", fmt.Sprintf("%+v", msc1))
	it.Equal("{Class:cls Message:msg Inner:<nil> StackTrace:\nstack}", fmt.Sprintf("%+v", msc2))

	// Format: *Ex -> +v
	it.Equal("", fmt.Sprintf("%+v", &e))
	it.Equal("cls", fmt.Sprintf("%+v", &c1))
	it.Equal("cls", fmt.Sprintf("%+v", &c2))
	it.Equal("; msg", fmt.Sprintf("%+v", &m))
	it.Equal("\ninner", fmt.Sprintf("%+v", &i1))
	it.Equal("\ninner", fmt.Sprintf("%+v", &i2))
	it.Equal("\n; inner", fmt.Sprintf("%+v", &i3))
	it.Equal("\ninner", fmt.Sprintf("%+v", &i4))
	it.Equal("\nstack", fmt.Sprintf("%+v", &s))
	it.Equal("cls; msg", fmt.Sprintf("%+v", &mc1))
	it.Equal("cls; msg", fmt.Sprintf("%+v", &mc2))
	it.Equal("cls; msg\nstack", fmt.Sprintf("%+v", &msc1))
	it.Equal("cls; msg\nstack", fmt.Sprintf("%+v", &msc2))

	// Format: Ex -> #v
	it.Equal("<nil>", fmt.Sprintf("%#v", n))
	it.Equal(`ex.Ex{Class:error(nil), Message:"", Inner:error(nil), StackTrace:ex.StackTrace(nil)}`, fmt.Sprintf("%#v", e))
	it.Equal(`ex.Ex{Class:"cls", Message:"", Inner:error(nil), StackTrace:ex.StackTrace(nil)}`, fmt.Sprintf("%#v", c1))
	it.HasPrefix(fmt.Sprintf("%#v", c2), `ex.Ex{Class:(*errors.errorString)(0x`)
	it.HasSuffix(fmt.Sprintf("%#v", c2), `), Message:"", Inner:error(nil), StackTrace:ex.StackTrace(nil)}`)
	it.Equal(`ex.Ex{Class:error(nil), Message:"msg", Inner:error(nil), StackTrace:ex.StackTrace(nil)}`, fmt.Sprintf("%#v", m))
	it.Equal(`ex.Ex{Class:error(nil), Message:"", Inner:"inner", StackTrace:ex.StackTrace(nil)}`, fmt.Sprintf("%#v", i1))
	it.Equal(`ex.Ex{Class:error(nil), Message:"", Inner:inner, StackTrace:ex.StackTrace(nil)}`, fmt.Sprintf("%#v", i2))
	it.Equal(`ex.Ex{Class:error(nil), Message:"", Inner:; inner, StackTrace:ex.StackTrace(nil)}`, fmt.Sprintf("%#v", i3))
	it.HasPrefix(fmt.Sprintf("%#v", i4), `ex.Ex{Class:error(nil), Message:"", Inner:(*errors.errorString)(0x`)
	it.HasSuffix(fmt.Sprintf("%#v", i4), `), StackTrace:ex.StackTrace(nil)}`)
	it.Equal(`ex.Ex{Class:error(nil), Message:"", Inner:error(nil), StackTrace:[]string{"stack"}}`, fmt.Sprintf("%#v", s))
	it.Equal(`ex.Ex{Class:"cls", Message:"msg", Inner:error(nil), StackTrace:ex.StackTrace(nil)}`, fmt.Sprintf("%#v", mc1))
	it.HasPrefix(fmt.Sprintf("%#v", mc2), `ex.Ex{Class:(*errors.errorString)(0x`)
	it.HasSuffix(fmt.Sprintf("%#v", mc2), `), Message:"msg", Inner:error(nil), StackTrace:ex.StackTrace(nil)}`)
	it.Equal(`ex.Ex{Class:"cls", Message:"msg", Inner:error(nil), StackTrace:[]string{"stack"}}`, fmt.Sprintf("%#v", msc1))
	it.HasPrefix(fmt.Sprintf("%#v", msc2), `ex.Ex{Class:(*errors.errorString)(0x`)
	it.HasSuffix(fmt.Sprintf("%#v", msc2), `), Message:"msg", Inner:error(nil), StackTrace:[]string{"stack"}}`)

	// Format: *Ex -> #v
	it.Equal("", fmt.Sprintf("%#v", &e))
	it.Equal("cls", fmt.Sprintf("%#v", &c1))
	it.Equal("cls", fmt.Sprintf("%#v", &c2))
	it.Equal("; msg", fmt.Sprintf("%#v", &m))
	it.Equal("\ninner", fmt.Sprintf("%#v", &i1))
	it.Equal("\ninner", fmt.Sprintf("%#v", &i2))
	it.Equal("\n; inner", fmt.Sprintf("%#v", &i3))
	it.Equal("\ninner", fmt.Sprintf("%#v", &i4))
	it.Equal("", fmt.Sprintf("%#v", &s))
	it.Equal("cls; msg", fmt.Sprintf("%#v", &mc1))
	it.Equal("cls; msg", fmt.Sprintf("%#v", &mc2))
	it.Equal("cls; msg", fmt.Sprintf("%#v", &msc1))
	it.Equal("cls; msg", fmt.Sprintf("%#v", &msc2))

	// Format: Ex -> c
	// NOTE: trick the linter
	var cfmt = "%"
	cfmt += "c"
	it.Equal("<nil>", fmt.Sprintf(cfmt, n))
	it.Equal("{<nil> %!c(string=) <nil> <nil>}", fmt.Sprintf(cfmt, e))
	it.Equal("{%!c(ex.Class=cls) %!c(string=) <nil> <nil>}", fmt.Sprintf(cfmt, c1))
	it.Equal("{%!c(*errors.errorString=&{cls}) %!c(string=) <nil> <nil>}", fmt.Sprintf(cfmt, c2))
	it.Equal("{<nil> %!c(string=msg) <nil> <nil>}", fmt.Sprintf(cfmt, m))
	it.Equal("{<nil> %!c(string=) %!c(ex.Class=inner) <nil>}", fmt.Sprintf(cfmt, i1))
	it.Equal("{<nil> %!c(string=) inner <nil>}", fmt.Sprintf(cfmt, i2))
	it.Equal("{<nil> %!c(string=) %!c(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference) <nil>}", fmt.Sprintf(cfmt, i3))
	it.Equal("{<nil> %!c(string=) %!c(*errors.errorString=&{inner}) <nil>}", fmt.Sprintf(cfmt, i4))
	it.Equal("{<nil> %!c(string=) <nil> }", fmt.Sprintf(cfmt, s))
	it.Equal("{%!c(ex.Class=cls) %!c(string=msg) <nil> <nil>}", fmt.Sprintf(cfmt, mc1))
	it.Equal("{%!c(*errors.errorString=&{cls}) %!c(string=msg) <nil> <nil>}", fmt.Sprintf(cfmt, mc2))
	it.Equal("{%!c(ex.Class=cls) %!c(string=msg) <nil> }", fmt.Sprintf(cfmt, msc1))
	it.Equal("{%!c(*errors.errorString=&{cls}) %!c(string=msg) <nil> }", fmt.Sprintf(cfmt, msc2))

	// Format: *Ex -> c
	it.Equal("%!c(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)", fmt.Sprintf("%c", &e))
	it.Equal("cls", fmt.Sprintf("%c", &c1))
	it.Equal("cls", fmt.Sprintf("%c", &c2))
	it.Equal("%!c(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)", fmt.Sprintf("%c", &m))
	it.Equal("%!c(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)", fmt.Sprintf("%c", &i1))
	it.Equal("%!c(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)", fmt.Sprintf("%c", &i2))
	it.Equal("%!c(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)", fmt.Sprintf("%c", &i3))
	it.Equal("%!c(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)", fmt.Sprintf("%c", &i4))
	it.Equal("%!c(PANIC=Format method: runtime error: invalid memory address or nil pointer dereference)", fmt.Sprintf("%c", &s))
	it.Equal("cls", fmt.Sprintf("%c", &mc1))
	it.Equal("cls", fmt.Sprintf("%c", &mc2))
	it.Equal("cls", fmt.Sprintf("%c", &msc1))
	it.Equal("cls", fmt.Sprintf("%c", &msc2))

	// Format: Ex -> i
	it.Equal("<nil>", fmt.Sprintf("%i", any(n)))
	it.Equal("{<nil> %!i(string=) <nil> <nil>}", fmt.Sprintf("%i", any(e)))
	it.Equal("{%!i(ex.Class=cls) %!i(string=) <nil> <nil>}", fmt.Sprintf("%i", any(c1)))
	it.Equal("{%!i(*errors.errorString=&{cls}) %!i(string=) <nil> <nil>}", fmt.Sprintf("%i", any(c2)))
	it.Equal("{<nil> %!i(string=msg) <nil> <nil>}", fmt.Sprintf("%i", any(m)))
	it.Equal("{<nil> %!i(string=) %!i(ex.Class=inner) <nil>}", fmt.Sprintf("%i", any(i1)))
	it.Equal("{<nil> %!i(string=)  <nil>}", fmt.Sprintf("%i", any(i2)))
	it.Equal("{<nil> %!i(string=)  <nil>}", fmt.Sprintf("%i", any(i3)))
	it.Equal("{<nil> %!i(string=) %!i(*errors.errorString=&{inner}) <nil>}", fmt.Sprintf("%i", any(i4)))
	it.Equal("{<nil> %!i(string=) <nil> }", fmt.Sprintf("%i", any(s)))
	it.Equal("{%!i(ex.Class=cls) %!i(string=msg) <nil> <nil>}", fmt.Sprintf("%i", any(mc1)))
	it.Equal("{%!i(*errors.errorString=&{cls}) %!i(string=msg) <nil> <nil>}", fmt.Sprintf("%i", any(mc2)))
	it.Equal("{%!i(ex.Class=cls) %!i(string=msg) <nil> }", fmt.Sprintf("%i", any(msc1)))
	it.Equal("{%!i(*errors.errorString=&{cls}) %!i(string=msg) <nil> }", fmt.Sprintf("%i", any(msc2)))

	// Format: *Ex -> i
	it.Equal("", fmt.Sprintf("%i", &e))
	it.Equal("", fmt.Sprintf("%i", &c1))
	it.Equal("", fmt.Sprintf("%i", &c2))
	it.Equal("", fmt.Sprintf("%i", &m))
	it.Equal("inner", fmt.Sprintf("%i", &i1))
	it.Equal("", fmt.Sprintf("%i", &i2))
	it.Equal("", fmt.Sprintf("%i", &i3))
	it.Equal("inner", fmt.Sprintf("%i", &i4))
	it.Equal("", fmt.Sprintf("%i", &s))
	it.Equal("", fmt.Sprintf("%i", &mc1))
	it.Equal("", fmt.Sprintf("%i", &mc2))
	it.Equal("", fmt.Sprintf("%i", &msc1))
	it.Equal("", fmt.Sprintf("%i", &msc2))

	// Format: Ex -> m
	it.Equal("<nil>", fmt.Sprintf("%m", any(n)))
	it.Equal("{<nil> %!m(string=) <nil> <nil>}", fmt.Sprintf("%m", any(e)))
	it.Equal("{%!m(ex.Class=cls) %!m(string=) <nil> <nil>}", fmt.Sprintf("%m", any(c1)))
	it.Equal("{%!m(*errors.errorString=&{cls}) %!m(string=) <nil> <nil>}", fmt.Sprintf("%m", any(c2)))
	it.Equal("{<nil> %!m(string=msg) <nil> <nil>}", fmt.Sprintf("%m", any(m)))
	it.Equal("{<nil> %!m(string=) %!m(ex.Class=inner) <nil>}", fmt.Sprintf("%m", any(i1)))
	it.Equal("{<nil> %!m(string=)  <nil>}", fmt.Sprintf("%m", any(i2)))
	it.Equal("{<nil> %!m(string=) inner <nil>}", fmt.Sprintf("%m", any(i3)))
	it.Equal("{<nil> %!m(string=) %!m(*errors.errorString=&{inner}) <nil>}", fmt.Sprintf("%m", any(i4)))
	it.Equal("{<nil> %!m(string=) <nil> }", fmt.Sprintf("%m", any(s)))
	it.Equal("{%!m(ex.Class=cls) %!m(string=msg) <nil> <nil>}", fmt.Sprintf("%m", any(mc1)))
	it.Equal("{%!m(*errors.errorString=&{cls}) %!m(string=msg) <nil> <nil>}", fmt.Sprintf("%m", any(mc2)))
	it.Equal("{%!m(ex.Class=cls) %!m(string=msg) <nil> }", fmt.Sprintf("%m", any(msc1)))
	it.Equal("{%!m(*errors.errorString=&{cls}) %!m(string=msg) <nil> }", fmt.Sprintf("%m", any(msc2)))

	// Format: *Ex -> m
	it.Equal("", fmt.Sprintf("%m", &e))
	it.Equal("", fmt.Sprintf("%m", &c1))
	it.Equal("", fmt.Sprintf("%m", &c2))
	it.Equal("msg", fmt.Sprintf("%m", &m))
	it.Equal("", fmt.Sprintf("%m", &i1))
	it.Equal("", fmt.Sprintf("%m", &i2))
	it.Equal("", fmt.Sprintf("%m", &i3))
	it.Equal("", fmt.Sprintf("%m", &i4))
	it.Equal("", fmt.Sprintf("%m", &s))
	it.Equal("msg", fmt.Sprintf("%m", &mc1))
	it.Equal("msg", fmt.Sprintf("%m", &mc2))
	it.Equal("msg", fmt.Sprintf("%m", &msc1))
	it.Equal("msg", fmt.Sprintf("%m", &msc2))

	// Format: Ex -> q
	it.Equal("<nil>", fmt.Sprintf("%q", n))
	it.Equal(`{<nil> "" <nil> <nil>}`, fmt.Sprintf("%q", e))
	it.Equal(`{"cls" "" <nil> <nil>}`, fmt.Sprintf("%q", c1))
	it.Equal(`{"cls" "" <nil> <nil>}`, fmt.Sprintf("%q", c2))
	it.Equal(`{<nil> "msg" <nil> <nil>}`, fmt.Sprintf("%q", m))
	it.Equal(`{<nil> "" "inner" <nil>}`, fmt.Sprintf("%q", i1))
	it.Equal(`{<nil> "" "" <nil>}`, fmt.Sprintf("%q", i2))
	it.Equal(`{<nil> "" "inner" <nil>}`, fmt.Sprintf("%q", i3))
	it.Equal(`{<nil> "" "inner" <nil>}`, fmt.Sprintf("%q", i4))
	it.Equal(`{<nil> "" <nil> }`, fmt.Sprintf("%q", s))
	it.Equal(`{"cls" "msg" <nil> <nil>}`, fmt.Sprintf("%q", mc1))
	it.Equal(`{"cls" "msg" <nil> <nil>}`, fmt.Sprintf("%q", mc2))
	it.Equal(`{"cls" "msg" <nil> }`, fmt.Sprintf("%q", msc1))
	it.Equal(`{"cls" "msg" <nil> }`, fmt.Sprintf("%q", msc2))

	// Format: *Ex -> q
	it.Equal(`""`, fmt.Sprintf("%q", &e))
	it.Equal(`""`, fmt.Sprintf("%q", &c1))
	it.Equal(`""`, fmt.Sprintf("%q", &c2))
	it.Equal(`"msg"`, fmt.Sprintf("%q", &m))
	it.Equal(`""`, fmt.Sprintf("%q", &i1))
	it.Equal(`""`, fmt.Sprintf("%q", &i2))
	it.Equal(`""`, fmt.Sprintf("%q", &i3))
	it.Equal(`""`, fmt.Sprintf("%q", &i4))
	it.Equal(`""`, fmt.Sprintf("%q", &s))
	it.Equal(`"msg"`, fmt.Sprintf("%q", &mc1))
	it.Equal(`"msg"`, fmt.Sprintf("%q", &mc2))
	it.Equal(`"msg"`, fmt.Sprintf("%q", &msc1))
	it.Equal(`"msg"`, fmt.Sprintf("%q", &msc2))

	// Format: Ex -> t
	// NOTE: trick the linter
	var tfmt = "%"
	tfmt += "t"
	it.Equal("", fmt.Sprintf(tfmt, n))
	it.Equal("{<nil> %!t(string=) <nil> <nil>}", fmt.Sprintf(tfmt, any(e)))
	it.Equal("{%!t(ex.Class=cls) %!t(string=) <nil> <nil>}", fmt.Sprintf(tfmt, any(c1)))
	it.Equal("{%!t(*errors.errorString=&{cls}) %!t(string=) <nil> <nil>}", fmt.Sprintf(tfmt, any(c2)))
	it.Equal("{<nil> %!t(string=msg) <nil> <nil>}", fmt.Sprintf(tfmt, any(m)))
	it.Equal("{<nil> %!t(string=) %!t(ex.Class=inner) <nil>}", fmt.Sprintf(tfmt, any(i1)))
	it.Equal("{<nil> %!t(string=)  <nil>}", fmt.Sprintf(tfmt, any(i2)))
	it.Equal("{<nil> %!t(string=)  <nil>}", fmt.Sprintf(tfmt, any(i3)))
	it.Equal("{<nil> %!t(string=) %!t(*errors.errorString=&{inner}) <nil>}", fmt.Sprintf(tfmt, any(i4)))
	it.Equal("{<nil> %!t(string=) <nil> }", fmt.Sprintf(tfmt, any(s)))
	it.Equal("{%!t(ex.Class=cls) %!t(string=msg) <nil> <nil>}", fmt.Sprintf(tfmt, any(mc1)))
	it.Equal("{%!t(*errors.errorString=&{cls}) %!t(string=msg) <nil> <nil>}", fmt.Sprintf(tfmt, any(mc2)))
	it.Equal("{%!t(ex.Class=cls) %!t(string=msg) <nil> }", fmt.Sprintf(tfmt, any(msc1)))
	it.Equal("{%!t(*errors.errorString=&{cls}) %!t(string=msg) <nil> }", fmt.Sprintf(tfmt, any(msc2)))

	// Format: *Ex -> t
	it.Equal("", fmt.Sprintf("%t", &e))
	it.Equal("", fmt.Sprintf("%t", &c1))
	it.Equal("", fmt.Sprintf("%t", &c2))
	it.Equal("", fmt.Sprintf("%t", &m))
	it.Equal("", fmt.Sprintf("%t", &i1))
	it.Equal("", fmt.Sprintf("%t", &i2))
	it.Equal("", fmt.Sprintf("%t", &i3))
	it.Equal("", fmt.Sprintf("%t", &i4))
	it.Equal("", fmt.Sprintf("%t", &s))
	it.Equal("", fmt.Sprintf("%t", &mc1))
	it.Equal("", fmt.Sprintf("%t", &mc2))
	it.Equal("", fmt.Sprintf("%t", &msc1))
	it.Equal("", fmt.Sprintf("%t", &msc2))

	// Format: Ex -> [default Sprint]
	it.Equal("<nil>", fmt.Sprint(n))
	it.Equal("{<nil>  <nil> <nil>}", fmt.Sprint(e))
	it.Equal("{cls  <nil> <nil>}", fmt.Sprint(c1))
	it.Equal("{cls  <nil> <nil>}", fmt.Sprint(c2))
	it.Equal("{<nil> msg <nil> <nil>}", fmt.Sprint(m))
	it.Equal("{<nil>  inner <nil>}", fmt.Sprint(i1))
	it.Equal("{<nil>  inner <nil>}", fmt.Sprint(i2))
	it.Equal("{<nil>  ; inner <nil>}", fmt.Sprint(i3))
	it.Equal("{<nil>  inner <nil>}", fmt.Sprint(i4))
	it.Equal("{<nil>  <nil> \nstack}", fmt.Sprint(s))
	it.Equal("{cls msg <nil> <nil>}", fmt.Sprint(mc1))
	it.Equal("{cls msg <nil> <nil>}", fmt.Sprint(mc2))
	it.Equal("{cls msg <nil> \nstack}", fmt.Sprint(msc1))
	it.Equal("{cls msg <nil> \nstack}", fmt.Sprint(msc2))

	// Format: *Ex -> [default Sprint]
	it.Equal("", fmt.Sprint(&e))
	it.Equal("cls", fmt.Sprint(&c1))
	it.Equal("cls", fmt.Sprint(&c2))
	it.Equal("; msg", fmt.Sprint(&m))
	it.Equal("\ninner", fmt.Sprint(&i1))
	it.Equal("\ninner", fmt.Sprint(&i2))
	it.Equal("\n; inner", fmt.Sprint(&i3))
	it.Equal("\ninner", fmt.Sprint(&i4))
	it.Equal("", fmt.Sprint(&s))
	it.Equal("cls; msg", fmt.Sprint(&mc1))
	it.Equal("cls; msg", fmt.Sprint(&mc2))
	it.Equal("cls; msg", fmt.Sprint(&msc1))
	it.Equal("cls; msg", fmt.Sprint(&msc2))
}

func TestMarshalJSON(t *testing.T) {

	type ReadableStackTrace struct {
		Class   string   `json:"Class"`
		Message string   `json:"Message"`
		Inner   error    `json:"Inner"`
		Stack   []string `json:"StackTrace"`
	}

	a := assert.New(t)
	message := "new test error"
	ex := As(New(message))
	a.NotNil(ex)
	stackTrace := ex.StackTrace
	typed, isTyped := stackTrace.(StackPointers)
	a.True(isTyped)
	a.NotNil(typed)
	stackDepth := len(typed)

	jsonErr, err := json.Marshal(ex)
	a.Nil(err)
	a.NotNil(jsonErr)

	ex2 := &ReadableStackTrace{}
	err = json.Unmarshal(jsonErr, ex2)
	a.Nil(err)
	a.Len(ex2.Stack, stackDepth)
	a.Equal(message, ex2.Class)

	ex = As(New(fmt.Errorf(message)))
	a.NotNil(ex)
	stackTrace = ex.StackTrace
	typed, isTyped = stackTrace.(StackPointers)
	a.True(isTyped)
	a.NotNil(typed)
	stackDepth = len(typed)

	jsonErr, err = json.Marshal(ex)
	a.Nil(err)
	a.NotNil(jsonErr)

	ex2 = &ReadableStackTrace{}
	err = json.Unmarshal(jsonErr, ex2)
	a.Nil(err)
	a.Len(ex2.Stack, stackDepth)
	a.Equal(message, ex2.Class)
}

func TestJSON(t *testing.T) {
	assert := assert.New(t)

	ex := New("this is a test",
		OptMessage("test message"),
		OptInner(New("inner exception", OptMessagef("inner test message"))),
	)

	contents, err := json.Marshal(ex)
	assert.Nil(err)

	var verify Ex
	err = json.Unmarshal(contents, &verify)
	assert.Nil(err)

	assert.Equal(ErrClass(ex), ErrClass(verify))
	assert.Equal(ErrMessage(ex), ErrMessage(verify))
	assert.NotNil(verify.Inner)
	assert.Equal(ErrClass(ErrInner(ex)), ErrClass(ErrInner(verify)))
	assert.Equal(ErrMessage(ErrInner(ex)), ErrMessage(ErrInner(verify)))
}

func TestNest(t *testing.T) {
	a := assert.New(t)

	ex1 := As(New("this is an error"))
	ex2 := As(New("this is another error"))
	err := As(Nest(ex1, ex2))

	a.NotNil(err)
	a.NotNil(err.Inner)
	a.NotEmpty(err.Error())

	a.True(Is(ex1, Class("this is an error")))
	a.True(Is(ex1.Inner, Class("this is another error")))
}

func TestNestNil(t *testing.T) {
	a := assert.New(t)

	var ex1 error
	var ex2 error
	var ex3 error

	err := Nest(ex1, ex2, ex3)
	a.Nil(err)
	a.Equal(nil, err)
	a.True(nil == err)
}

func TestExceptionFormat(t *testing.T) {
	assert := assert.New(t)

	e := &Ex{Class: fmt.Errorf("this is only a test")}
	output := fmt.Sprintf("%v", e)
	assert.Equal("this is only a test", output)

	output = fmt.Sprintf("%+v", e)
	assert.Equal("this is only a test", output)

	e = &Ex{
		Class: fmt.Errorf("this is only a test"),
		StackTrace: StackStrings([]string{
			"foo",
			"bar",
		}),
	}

	output = fmt.Sprintf("%+v", e)
	assert.Equal("this is only a test\nfoo\nbar", output)
}

func TestExceptionPrintsInner(t *testing.T) {
	assert := assert.New(t)

	ex := New("outer", OptInner(New("middle", OptInner(New("terminal")))))

	output := fmt.Sprintf("%v", ex)

	assert.Contains(output, "outer")
	assert.Contains(output, "middle")
	assert.Contains(output, "terminal")

	output = fmt.Sprintf("%+v", ex)

	assert.Contains(output, "outer")
	assert.Contains(output, "middle")
	assert.Contains(output, "terminal")
}

type structuredError struct {
	value string
}

func (err structuredError) Error() string {
	return err.value
}

func TestException_ErrorsIsCompatability(t *testing.T) {
	assert := assert.New(t)

	{ // Single nesting, Ex is outermost
		innerErr := errors.New("inner")
		outerErr := New("outer", OptInnerClass(innerErr))

		assert.True(errors.Is(outerErr, innerErr))
	}

	{ // Single nesting, Ex is innermost
		innerErr := New("inner")
		outerErr := fmt.Errorf("outer: %w", innerErr)

		assert.True(errors.Is(outerErr, Class("inner")))
	}

	{ // Triple nesting, including Ex and non-Ex
		firstErr := errors.New("inner most")
		secondErr := fmt.Errorf("standard err: %w", firstErr)
		thirdErr := New("ex err", OptInner(secondErr))
		fourthErr := New("outer most", OptInner(thirdErr))

		assert.True(errors.Is(fourthErr, firstErr))
		assert.True(errors.Is(fourthErr, secondErr))
		assert.True(errors.Is(fourthErr, Class("ex err")))
	}

	{ // Target is nested in an Ex class and not in Inner chain
		firstErr := errors.New("inner most")
		secondErr := fmt.Errorf("standard err: %w", firstErr)
		thirdErr := New(secondErr, OptInner(fmt.Errorf("another cause")))

		assert.True(errors.Is(thirdErr, firstErr))
		assert.True(errors.Is(thirdErr, secondErr))
	}

	{ // Simple Ex against a Multi
		firstErr := New("ex")
		secondErr := errors.New("second")
		thirdErr := Multi{firstErr, secondErr}

		assert.False(errors.Is(firstErr, thirdErr))
	}

	{ // Ex wrapping a Multi testing against Multi
		firstErr := errors.New("inner most")
		secondErr := fmt.Errorf("standard err: %w", firstErr)
		thirdErr := New(Multi{firstErr, secondErr})

		assert.True(errors.Is(thirdErr, firstErr))
		assert.True(errors.Is(thirdErr, secondErr))
		assert.True(errors.Is(thirdErr, thirdErr))
		assert.False(errors.Is(thirdErr, Multi{firstErr, secondErr}))
		assert.False(errors.Is(thirdErr, Multi{firstErr, secondErr, secondErr}))
		assert.False(errors.Is(thirdErr, Multi{secondErr, firstErr}))
	}

	{
		// nil checks
		var nilEx *Ex
		assert.False(errors.Is(nilEx, nil))
		assert.False(errors.Is(&Ex{}, nil))
	}
}

func TestException_ErrorsAsCompatability(t *testing.T) {
	assert := assert.New(t)

	{ // Single nesting, targeting non-Ex
		innerErr := structuredError{"inner most"}
		outerErr := New("outer", OptInner(innerErr))

		var matchedErr structuredError
		assert.True(errors.As(outerErr, &matchedErr))
		assert.Equal("inner most", matchedErr.value)
	}

	{ // Single nesting, targeting Ex
		innerErr := New("outer most")
		outerErr := fmt.Errorf("outer err: %w", innerErr)

		var matchedErr *Ex
		assert.True(errors.As(outerErr, &matchedErr))
		assert.Equal("outer most", matchedErr.Class.Error())
	}

	{ // Single nesting, targeting inner Ex class
		innerErr := New(structuredError{"inner most"})
		outerErr := New("outer most", OptInner(innerErr))

		var matchedErr structuredError
		assert.True(errors.As(outerErr, &matchedErr))
		assert.Equal("inner most", matchedErr.value)
	}

	{ // Triple Nesting, targeting non-Ex
		firstErr := structuredError{"inner most"}
		secondErr := fmt.Errorf("standard err: %w", firstErr)
		thirdErr := New("ex err", OptInner(secondErr))
		fourthErr := New("outer most", OptInner(thirdErr))

		var matchedErr structuredError
		assert.True(errors.As(fourthErr, &matchedErr))
		assert.Equal("inner most", matchedErr.value)
	}
}
