/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package web

import (
	"net/http"
	"testing"
	"time"

	"github.com/blend/go-sdk/assert"
)

var readSetCookiesTests = []struct {
	Header	http.Header
	Cookies	[]*http.Cookie
}{
	{
		http.Header{"Set-Cookie": {"Cookie-1=v$1"}},
		[]*http.Cookie{{Name: "Cookie-1", Value: "v$1", Raw: "Cookie-1=v$1"}},
	},
	{
		http.Header{"Set-Cookie": {"NID=99=YsDT5i3E-CXax-; expires=Wed, 23-Nov-2011 01:05:03 GMT; path=/; domain=.google.ch; HttpOnly"}},
		[]*http.Cookie{{
			Name:		"NID",
			Value:		"99=YsDT5i3E-CXax-",
			Path:		"/",
			Domain:		".google.ch",
			HttpOnly:	true,
			Expires:	time.Date(2011, 11, 23, 1, 5, 3, 0, time.UTC),
			RawExpires:	"Wed, 23-Nov-2011 01:05:03 GMT",
			Raw:		"NID=99=YsDT5i3E-CXax-; expires=Wed, 23-Nov-2011 01:05:03 GMT; path=/; domain=.google.ch; HttpOnly",
		}},
	},
	{
		http.Header{"Set-Cookie": {".ASPXAUTH=7E3AA; expires=Wed, 07-Mar-2012 14:25:06 GMT; path=/; HttpOnly"}},
		[]*http.Cookie{{
			Name:		".ASPXAUTH",
			Value:		"7E3AA",
			Path:		"/",
			Expires:	time.Date(2012, 3, 7, 14, 25, 6, 0, time.UTC),
			RawExpires:	"Wed, 07-Mar-2012 14:25:06 GMT",
			HttpOnly:	true,
			Raw:		".ASPXAUTH=7E3AA; expires=Wed, 07-Mar-2012 14:25:06 GMT; path=/; HttpOnly",
		}},
	},
	{
		http.Header{"Set-Cookie": {"ASP.NET_SessionId=foo; path=/; HttpOnly"}},
		[]*http.Cookie{{
			Name:		"ASP.NET_SessionId",
			Value:		"foo",
			Path:		"/",
			HttpOnly:	true,
			Raw:		"ASP.NET_SessionId=foo; path=/; HttpOnly",
		}},
	},
	{
		http.Header{"Set-Cookie": {"samesitedefault=foo; SameSite"}},
		[]*http.Cookie{{
			Name:		"samesitedefault",
			Value:		"foo",
			SameSite:	http.SameSiteDefaultMode,
			Raw:		"samesitedefault=foo; SameSite",
		}},
	},
	{
		http.Header{"Set-Cookie": {"samesitelax=foo; SameSite=Lax"}},
		[]*http.Cookie{{
			Name:		"samesitelax",
			Value:		"foo",
			SameSite:	http.SameSiteLaxMode,
			Raw:		"samesitelax=foo; SameSite=Lax",
		}},
	},
	{
		http.Header{"Set-Cookie": {"samesitestrict=foo; SameSite=Strict"}},
		[]*http.Cookie{{
			Name:		"samesitestrict",
			Value:		"foo",
			SameSite:	http.SameSiteStrictMode,
			Raw:		"samesitestrict=foo; SameSite=Strict",
		}},
	},
	// Make sure we can properly read back the Set-Cookie headers we create
	// for values containing spaces or commas:
	{
		http.Header{"Set-Cookie": {`special-1=a z`}},
		[]*http.Cookie{{Name: "special-1", Value: "a z", Raw: `special-1=a z`}},
	},
	{
		http.Header{"Set-Cookie": {`special-2=" z"`}},
		[]*http.Cookie{{Name: "special-2", Value: " z", Raw: `special-2=" z"`}},
	},
	{
		http.Header{"Set-Cookie": {`special-3="a "`}},
		[]*http.Cookie{{Name: "special-3", Value: "a ", Raw: `special-3="a "`}},
	},
	{
		http.Header{"Set-Cookie": {`special-4=" "`}},
		[]*http.Cookie{{Name: "special-4", Value: " ", Raw: `special-4=" "`}},
	},
	{
		http.Header{"Set-Cookie": {`special-5=a,z`}},
		[]*http.Cookie{{Name: "special-5", Value: "a,z", Raw: `special-5=a,z`}},
	},
	{
		http.Header{"Set-Cookie": {`special-6=",z"`}},
		[]*http.Cookie{{Name: "special-6", Value: ",z", Raw: `special-6=",z"`}},
	},
	{
		http.Header{"Set-Cookie": {`special-7=a,`}},
		[]*http.Cookie{{Name: "special-7", Value: "a,", Raw: `special-7=a,`}},
	},
	{
		http.Header{"Set-Cookie": {`special-8=","`}},
		[]*http.Cookie{{Name: "special-8", Value: ",", Raw: `special-8=","`}},
	},

	// TODO(bradfitz): users have reported seeing this in the
	// wild, but do browsers handle it? RFC 6265 just says "don't
	// do that" (section 3) and then never mentions header folding
	// again.
	// Header{"Set-Cookie": {"ASP.NET_SessionId=foo; path=/; HttpOnly, .ASPXAUTH=7E3AA; expires=Wed, 07-Mar-2012 14:25:06 GMT; path=/; HttpOnly"}},
}

func TestReadSetCookies(t *testing.T) {
	assert := assert.New(t)

	for _, tt := range readSetCookiesTests {
		for n := 0; n < 2; n++ {	// to verify readSetCookies doesn't mutate its input
			c := ReadSetCookies(tt.Header)
			assert.NonFatal().Equal(c, tt.Cookies)
		}
	}
}
