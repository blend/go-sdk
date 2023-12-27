/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

*/

package envoyutil_test

import (
	"crypto/x509"
	"fmt"
	"net/url"
	"regexp"
	"testing"

	sdkAssert "github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/certutil"
	"github.com/blend/go-sdk/ex"

	"github.com/blend/go-sdk/envoyutil"
)

const (
	fullXFCCTest				= `By=spiffe://cluster.local/ns/blend/sa/yule;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client";URI=spiffe://cluster.local/ns/blend/sa/cheer`
	xfccElementByTest			= `By=spiffe://cluster.local/ns/blend/sa/tide`
	xfccElementHashTest			= `Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688`
	xfccElementCertTest			= `Cert=` + xfccElementTestCertEncoded
	xfccElementChainTest			= `CHAIN=` + xfccElementTestCertEncoded
	xfccElementSubjectTest			= `SUBJECT="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client"`
	xfccElementURITest			= `URI=spiffe://cluster.local/ns/blend/sa/quasar`
	xfccElementDNSTest			= `dns=http://frontend.lyft.com`
	xfccElementEndTest			= `dns=http://frontend.lyft.com;`
	xfccElementNoneTest			= `key=value;dns=http://frontend.lyft.com`
	xfccElementMultiTest			= `By=spiffe://cluster.local/ns/blend/sa/laser;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688`
	xfccElementMalformedKeyTest		= `=value`
	xfccElementMultiMalformedKeyTest	= `=value;dns=http://frontend.lyft.com`
	xfccElementMultiCertTest		= `cert=` + xfccElementTestCertEncoded + xfccElementTestCert

	xfccElementMalformedEncoding	= "%"

	xfccElementTestCertEncoded	= `-----BEGIN%20CERTIFICATE-----%0AMIIFKjCCAxICCQCA5%2FOCxg%2FqiDANBgkqhkiG9w0BAQsFADBXMQswCQYDVQQGEwJV%0AUzELMAkGA1UECAwCQ0ExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoM%0ABEx5ZnQxFDASBgNVBAMMC1Rlc3QgQ2xpZW50MB4XDTIwMDYwNDE3NDkzNVoXDTIx%0AMDYwNDE3NDkzNVowVzELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRYwFAYDVQQH%0ADA1TYW4gRnJhbmNpc2NvMQ0wCwYDVQQKDARMeWZ0MRQwEgYDVQQDDAtUZXN0IENs%0AaWVudDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAKs6T8vcb8rIIkC4%0Aiz9h%2FOj6Iv%2BfazTLwNLK%2Fk58Ape5ZL0IdW6h8pWDlGnGz4X%2FTaJ5TwlamFo1h62v%0AsR8HPNOoLY0wmC2qHVquPF6eR9Lt5ejJiakr%2BYvf%2BU6LHXlOpOoot5rcTGoGBCf0%0AH3zmjdOE0o6hwJxMf54XQEVwNXqRrIDbY27mYS8eAVcSMrPUQVZ%2B3Vk1S56Imybz%0Adegi79IIoc6TzE5M7ChfJZBNNNZT08haJe6Oi%2FIgZhK3IexssY%2BQyD5uBSc7Mpas%0A6TstzeevIbeFy3Od2GhUy2Hz98qW%2FoO5iuerEArkNs4lB0J%2F0ARPHUDnmmH%2BqWYF%0APKealq2yEyXHHXrhDcSK%2FN5R64pp%2FVrxEas1qG20%2FCG4rixv36UJuEz5oUKNWyaR%0A268EI5Vecw%2BpK%2F0XC2%2Bhra9T%2FeP9JH0Fp43x7bdpQoxph8ZJZBsjbgCFMonf3ku1%0A9n74%2FxwvV6B0wp5C8jpwbGa85n%2BT8hogtO78mnpvxhTVJ7TOy596tI2apJ02edtD%0AJgsJV9MfZ%2FfGu3QZ6yN3rKVMPkZfC18cK04xy%2BroPo756CHkUHP5cz%2BKtJ7%2B8COR%0ArPDPxKBLOqwaSFcanQNONFIrffnZciiisCxjMHGoM4%2Fuix5gStlDC9%2FM5yyHt9He%0AldC8xL%2FyIalsa9Df7SL59Fd7T2JrAgMBAAEwDQYJKoZIhvcNAQELBQADggIBAGTb%0AOTddb6Yr37JigDYGKxSiidWfPgxTFr3DWTTE2aBhTi8K1Gr8Br1FAGRg4P114kxm%0AQBx3TxuCZcssMX2Z3GbL20DrHNExCl7hM%2FZVA2OCVhwfXGOHRQ6lcpmeQISWDNsN%0Atanlap%2FAgqKN%2F6At6JEYmuTSJnKc4Bfgk2GP5LPa63yJOlyvFb8ovKsCgb1ppVyw%0ARE%2B7AmB2DfDdVql4nHsDh5UBZRgVxMZ6xGnkYKaAUDKl4slejvKwXuzu2Xf%2BAd74%0AgjdLHzP0WmHlAggR5LIv%2F9xlvrsKCrNDDxWwOGeYk2WZl%2Fybud0RFKhLIqbbeMy7%0ADcdy04cJcqa9qRHYySgaWtM6Ab%2Fx9CJqdzR2NQZNnLgk6Vc3%2BoDjXMUuyM17WJAS%0ArenwJvanXvF9P1yPMByJQlXxkUehkCa%2FPs7E1O%2F%2BE2FJnvrtGVdYVR8Otbec1osS%0AmtJC6k7rgMhgvk63sCqQqaZwRWwLl2R5XcDZknUiqDKjuVHHA01II7jtGB1oyEIH%0Asp%2FrQlLNeyYlyhAlc3MhF5hu6nUjH%2B2%2BDuIHJsM0mEF0rjlbnp4bKJ%2FgF1COAIAL%0APzu2qAC%2BaOFldCmRonqUluayv6fQaQCeeh8sW2IjNVjA2ynKn2ybGIXH4mrH0KVa%0AJmUY%2B1YGMn7qbeHTma33N28Ec7hK%2BWByul746Nro%0A-----END%20CERTIFICATE-----`
	xfccElementTestCert		= `-----BEGIN CERTIFICATE-----
MIIFKjCCAxICCQCA5/OCxg/qiDANBgkqhkiG9w0BAQsFADBXMQswCQYDVQQGEwJV
UzELMAkGA1UECAwCQ0ExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoM
BEx5ZnQxFDASBgNVBAMMC1Rlc3QgQ2xpZW50MB4XDTIwMDYwNDE3NDkzNVoXDTIx
MDYwNDE3NDkzNVowVzELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRYwFAYDVQQH
DA1TYW4gRnJhbmNpc2NvMQ0wCwYDVQQKDARMeWZ0MRQwEgYDVQQDDAtUZXN0IENs
aWVudDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAKs6T8vcb8rIIkC4
iz9h/Oj6Iv+fazTLwNLK/k58Ape5ZL0IdW6h8pWDlGnGz4X/TaJ5TwlamFo1h62v
sR8HPNOoLY0wmC2qHVquPF6eR9Lt5ejJiakr+Yvf+U6LHXlOpOoot5rcTGoGBCf0
H3zmjdOE0o6hwJxMf54XQEVwNXqRrIDbY27mYS8eAVcSMrPUQVZ+3Vk1S56Imybz
degi79IIoc6TzE5M7ChfJZBNNNZT08haJe6Oi/IgZhK3IexssY+QyD5uBSc7Mpas
6TstzeevIbeFy3Od2GhUy2Hz98qW/oO5iuerEArkNs4lB0J/0ARPHUDnmmH+qWYF
PKealq2yEyXHHXrhDcSK/N5R64pp/VrxEas1qG20/CG4rixv36UJuEz5oUKNWyaR
268EI5Vecw+pK/0XC2+hra9T/eP9JH0Fp43x7bdpQoxph8ZJZBsjbgCFMonf3ku1
9n74/xwvV6B0wp5C8jpwbGa85n+T8hogtO78mnpvxhTVJ7TOy596tI2apJ02edtD
JgsJV9MfZ/fGu3QZ6yN3rKVMPkZfC18cK04xy+roPo756CHkUHP5cz+KtJ7+8COR
rPDPxKBLOqwaSFcanQNONFIrffnZciiisCxjMHGoM4/uix5gStlDC9/M5yyHt9He
ldC8xL/yIalsa9Df7SL59Fd7T2JrAgMBAAEwDQYJKoZIhvcNAQELBQADggIBAGTb
OTddb6Yr37JigDYGKxSiidWfPgxTFr3DWTTE2aBhTi8K1Gr8Br1FAGRg4P114kxm
QBx3TxuCZcssMX2Z3GbL20DrHNExCl7hM/ZVA2OCVhwfXGOHRQ6lcpmeQISWDNsN
tanlap/AgqKN/6At6JEYmuTSJnKc4Bfgk2GP5LPa63yJOlyvFb8ovKsCgb1ppVyw
RE+7AmB2DfDdVql4nHsDh5UBZRgVxMZ6xGnkYKaAUDKl4slejvKwXuzu2Xf+Ad74
gjdLHzP0WmHlAggR5LIv/9xlvrsKCrNDDxWwOGeYk2WZl/ybud0RFKhLIqbbeMy7
Dcdy04cJcqa9qRHYySgaWtM6Ab/x9CJqdzR2NQZNnLgk6Vc3+oDjXMUuyM17WJAS
renwJvanXvF9P1yPMByJQlXxkUehkCa/Ps7E1O/+E2FJnvrtGVdYVR8Otbec1osS
mtJC6k7rgMhgvk63sCqQqaZwRWwLl2R5XcDZknUiqDKjuVHHA01II7jtGB1oyEIH
sp/rQlLNeyYlyhAlc3MhF5hu6nUjH+2+DuIHJsM0mEF0rjlbnp4bKJ/gF1COAIAL
Pzu2qAC+aOFldCmRonqUluayv6fQaQCeeh8sW2IjNVjA2ynKn2ybGIXH4mrH0KVa
JmUY+1YGMn7qbeHTma33N28Ec7hK+WByul746Nro
-----END CERTIFICATE-----`
)

func TestXFCCElementDecodeBy(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		By		string
		Expected	*url.URL
		ErrorRegexp	*regexp.Regexp
	}
	testCases := []testCase{
		{By: "", Expected: &url.URL{}},
		// NOTE: Regex needed to support error format changes from go1.13 to go1.14
		{By: "\n", ErrorRegexp: regexp.MustCompile(`(?m)^parse ("\\n"|\n): net/url: invalid control character in URL$`)},
		{
			By:	"spiffe://cluster.local/ns/blend/sa/yule",
			Expected: &url.URL{
				Scheme:	"spiffe",
				Host:	"cluster.local",
				Path:	"/ns/blend/sa/yule",
			},
		},
	}
	for _, tc := range testCases {
		xe := envoyutil.XFCCElement{By: tc.By}
		uri, err := xe.DecodeBy()
		assert.Equal(tc.Expected, uri)
		if tc.ErrorRegexp != nil {
			asEx, ok := err.(*ex.Ex)
			assert.True(ok)
			assert.Equal(envoyutil.ErrXFCCParsing, asEx.Class)
			assert.True(tc.ErrorRegexp.MatchString(asEx.Inner.Error()))
		} else {
			assert.Nil(err)
		}
	}
}

func TestXFCCElementDecodeHash(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		Hash		string
		Expected	[]byte
		Error		string
	}
	testCases := []testCase{
		{Hash: "", Expected: []byte("")},
		{Hash: "41434944", Expected: []byte("ACID")},
		{Hash: "x", Expected: nil, Error: "Error Parsing X-Forwarded-Client-Cert\nencoding/hex: invalid byte: U+0078 'x'"},
	}
	for _, tc := range testCases {
		xe := envoyutil.XFCCElement{Hash: tc.Hash}
		decoded, err := xe.DecodeHash()
		assert.Equal(tc.Expected, decoded)
		if tc.Error != "" {
			assert.Equal(tc.Error, fmt.Sprintf("%v", err))
		} else {
			assert.Nil(err)
		}
	}
}

func TestXFCCElementDecodeCert(t *testing.T) {
	assert := sdkAssert.New(t)

	parsedCert, err := certutil.ParseCertPEM([]byte(xfccElementTestCert))
	assert.Nil(err)

	type testCase struct {
		Cert	string
		Parsed	*x509.Certificate
		Error	string
	}
	testCases := []testCase{
		{Cert: ""},
		{Cert: xfccElementTestCertEncoded, Parsed: parsedCert[0]},
		{Cert: "%", Error: "Error Parsing X-Forwarded-Client-Cert\ninvalid URL escape \"%\""},
		{
			Cert:	"-----BEGIN CERTIFICATE-----\nnope\n-----END CERTIFICATE-----\n",
			Error:	"Error Parsing X-Forwarded-Client-Cert\nx509: malformed certificate",
		},
		{
			Cert:	url.QueryEscape(xfccElementTestCert + "\n" + xfccElementTestCert),
			Error:	"Error Parsing X-Forwarded-Client-Cert; Incorrect number of certificates; expected 1 got 2",
		},
		{
			Cert:	xfccElementMultiCertTest,
			Error:	"Error Parsing X-Forwarded-Client-Cert; Incorrect number of certificates; expected 1 got 0",
		},
	}

	for _, tc := range testCases {
		xe := envoyutil.XFCCElement{Cert: tc.Cert}
		cert, err := xe.DecodeCert()
		if tc.Error != "" {
			assert.Equal(tc.Error, fmt.Sprintf("%v", err))
		} else {
			assert.Nil(err)
		}
		assert.Equal(tc.Parsed, cert)
	}
}

func TestXFCCElementDecodeChain(t *testing.T) {
	assert := sdkAssert.New(t)

	parsedCerts, err := certutil.ParseCertPEM([]byte(xfccElementTestCert + "\n" + xfccElementTestCert))
	assert.Nil(err)

	type testCase struct {
		Chain	string
		Parsed	[]*x509.Certificate
		Error	string
	}
	testCases := []testCase{
		{Chain: ""},
		{
			Chain:	url.QueryEscape(xfccElementTestCert + "\n" + xfccElementTestCert),
			Parsed:	parsedCerts,
		},
		{Chain: "%", Error: "Error Parsing X-Forwarded-Client-Cert\ninvalid URL escape \"%\""},
		{
			Chain:	"-----BEGIN CERTIFICATE-----\nnope\n-----END CERTIFICATE-----\n",
			Error:	"Error Parsing X-Forwarded-Client-Cert\nx509: malformed certificate",
		},
	}

	for _, tc := range testCases {
		xe := envoyutil.XFCCElement{Chain: tc.Chain}
		chain, err := xe.DecodeChain()
		if tc.Error != "" {
			assert.Equal(tc.Error, fmt.Sprintf("%v", err))
		} else {
			assert.Nil(err)
		}
		assert.Equal(tc.Parsed, chain)
	}
}

func TestXFCCElementDecodeURI(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		URI		string
		Expected	*url.URL
		ErrorRegexp	*regexp.Regexp
	}
	testCases := []testCase{
		{URI: "", Expected: &url.URL{}},
		// NOTE: Regex needed to support error format changes from go1.13 to go1.14
		{URI: "\r", ErrorRegexp: regexp.MustCompile(`(?m)^parse ("\\r"|\r): net/url: invalid control character in URL$`)},
		{
			URI:	"spiffe://cluster.local/ns/first/sa/furst",
			Expected: &url.URL{
				Scheme:	"spiffe",
				Host:	"cluster.local",
				Path:	"/ns/first/sa/furst",
			},
		},
	}
	for _, tc := range testCases {
		xe := envoyutil.XFCCElement{URI: tc.URI}
		uri, err := xe.DecodeURI()
		assert.Equal(tc.Expected, uri)
		if tc.ErrorRegexp != nil {
			asEx, ok := err.(*ex.Ex)
			assert.True(ok)
			assert.Equal(envoyutil.ErrXFCCParsing, asEx.Class)
			assert.True(tc.ErrorRegexp.MatchString(asEx.Inner.Error()))
		} else {
			assert.Nil(err)
		}
	}
}

func TestXFCCElementString(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		Element		envoyutil.XFCCElement
		Expected	string
	}
	testCases := []testCase{
		{Element: envoyutil.XFCCElement{By: "hi"}, Expected: "By=hi"},
		{Element: envoyutil.XFCCElement{Hash: "1389ab1"}, Expected: "Hash=1389ab1"},
		{Element: envoyutil.XFCCElement{Cert: "anything-goes"}, Expected: "Cert=anything-goes"},
		{Element: envoyutil.XFCCElement{Chain: "anything-goes"}, Expected: "Chain=anything-goes"},
		{
			Element:	envoyutil.XFCCElement{Subject: "OU=Blent/CN=Test Client"},
			Expected:	`Subject="OU=Blent/CN=Test Client"`,
		},
		{Element: envoyutil.XFCCElement{URI: "bye"}, Expected: "URI=bye"},
		{
			Element:	envoyutil.XFCCElement{DNS: []string{"web.invalid", "bye.invalid"}},
			Expected:	"DNS=web.invalid;DNS=bye.invalid",
		},
		{
			Element:	envoyutil.XFCCElement{By: "a,b=10", URI: `c; "then" again`},
			Expected:	"By=\"a,b=10\";URI=\"c; \\\"then\\\" again\"",
		},
	}
	for _, tc := range testCases {
		asString := tc.Element.String()
		assert.Equal(tc.Expected, asString)
		parsed, err := envoyutil.ParseXFCC(asString)
		assert.Nil(err)
		assert.Equal(envoyutil.XFCC{tc.Element}, parsed)
	}

	element := envoyutil.XFCCElement{}
	asString := element.String()
	assert.Equal("", asString)
}

func TestParseXFCC(t *testing.T) {
	assert := sdkAssert.New(t)

	ele, err := envoyutil.ParseXFCC(xfccElementByTest)
	assert.Nil(err)
	expected := envoyutil.XFCC{
		envoyutil.XFCCElement{By: "spiffe://cluster.local/ns/blend/sa/tide"},
	}
	assert.Equal(expected, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementByTest + ";" + xfccElementByTest)
	except, ok := err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr := &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	`Key already encountered "by"`,
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementHashTest)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			Hash: "468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688",
		},
	}
	assert.Equal(expected, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementHashTest + ";" + xfccElementHashTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	`Key already encountered "hash"`,
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementCertTest)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			Cert: xfccElementTestCertEncoded,
		},
	}
	assert.Equal(expected, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementCertTest + ";" + xfccElementCertTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	`Key already encountered "cert"`,
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementChainTest)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			Chain: xfccElementTestCertEncoded,
		},
	}
	assert.Equal(expected, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementChainTest + ";" + xfccElementChainTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	`Key already encountered "chain"`,
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementSubjectTest)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			Subject: "/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client",
		},
	}
	assert.Equal(expected, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementSubjectTest + ";" + xfccElementSubjectTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	`Key already encountered "subject"`,
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementURITest)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			URI: "spiffe://cluster.local/ns/blend/sa/quasar",
		},
	}
	assert.Equal(expected, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementURITest + ";" + xfccElementURITest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	`Key already encountered "uri"`,
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	ele, err = envoyutil.ParseXFCC(xfccElementDNSTest)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			DNS: []string{"http://frontend.lyft.com"},
		},
	}
	assert.Equal(expected, ele)

	ele, err = envoyutil.ParseXFCC("dns=web.invalid;dns=blend.local.invalid")
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			DNS: []string{"web.invalid", "blend.local.invalid"},
		},
	}
	assert.Equal(expected, ele)

	_, err = envoyutil.ParseXFCC(xfccElementNoneTest)
	assert.NotNil(err)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)

	ele, err = envoyutil.ParseXFCC(xfccElementMultiTest)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			By:	"spiffe://cluster.local/ns/blend/sa/laser",
			Hash:	"468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688",
		},
	}
	assert.Equal(expected, ele)

	_, err = envoyutil.ParseXFCC(xfccElementMalformedKeyTest)
	assert.NotNil(err)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)

	_, err = envoyutil.ParseXFCC(xfccElementMultiMalformedKeyTest)
	assert.NotNil(err)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)

	ele, err = envoyutil.ParseXFCC(xfccElementEndTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	"Ends with separator character",
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	ele, err = envoyutil.ParseXFCC("cert=" + xfccElementMalformedEncoding)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			Cert: "%",
		},
	}
	assert.Equal(expected, ele)

	ele, err = envoyutil.ParseXFCC("chain=" + xfccElementMalformedEncoding)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			Chain: "%",
		},
	}
	assert.Equal(expected, ele)

	ele, err = envoyutil.ParseXFCC("=;")
	assert.NotNil(err)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)
	assert.Equal(envoyutil.XFCC{}, ele)

	// Test empty subject
	ele, err = envoyutil.ParseXFCC(`By=spiffe://cluster.local/ns/blend/sa/protocol;Hash=52114972613efb0820c5e32bfee0f0ee2a84859f7169da6c222300ef852a1129;Subject="";URI=spiffe://cluster.local/ns/blend/sa/world`)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			By:		"spiffe://cluster.local/ns/blend/sa/protocol",
			Hash:		"52114972613efb0820c5e32bfee0f0ee2a84859f7169da6c222300ef852a1129",
			Subject:	"",
			URI:		"spiffe://cluster.local/ns/blend/sa/world",
		},
	}
	assert.Equal(expected, ele)

	// Quoted value with empty key.
	ele, err = envoyutil.ParseXFCC(`="a";b=20`)
	assert.Equal(envoyutil.XFCC{}, ele)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	"Key missing",
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)

	// Quoted value with invalid key.
	ele, err = envoyutil.ParseXFCC(`wrong="quoted";by=next`)
	assert.Equal(envoyutil.XFCC{}, ele)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	`Unknown key "wrong"`,
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)

	// Closing quoted not following by `;`
	ele, err = envoyutil.ParseXFCC(`a="b"---`)
	assert.Equal(envoyutil.XFCC{}, ele)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	"Closing quote not followed by `;`.",
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)

	// Escaped quotes and other characters work as expected.
	ele, err = envoyutil.ParseXFCC(`By="a,b=10";URI="c; \"then\" again"`)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			By:	"a,b=10",
			URI:	`c; "then" again`,
		},
	}
	assert.Equal(expected, ele)

	// Bare escape character works fine (when not followed by a quote).
	ele, err = envoyutil.ParseXFCC(`By="first\tsecond"`)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			By: "first\\tsecond",
		},
	}
	assert.Equal(expected, ele)

	xfcc, err := envoyutil.ParseXFCC(fullXFCCTest + "," + fullXFCCTest)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{
			By:		"spiffe://cluster.local/ns/blend/sa/yule",
			Hash:		"468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688",
			Subject:	"/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client",
			URI:		"spiffe://cluster.local/ns/blend/sa/cheer",
		},
		envoyutil.XFCCElement{
			By:		"spiffe://cluster.local/ns/blend/sa/yule",
			Hash:		"468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688",
			Subject:	"/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client",
			URI:		"spiffe://cluster.local/ns/blend/sa/cheer",
		},
	}
	assert.Equal(expected, xfcc)

	ele, err = envoyutil.ParseXFCC(xfccElementMalformedKeyTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	"Key or value found but not both",
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	// Quoted value at element boundary.
	xfcc, err = envoyutil.ParseXFCC(`by="me",uri=you`)
	assert.Nil(err)
	expected = envoyutil.XFCC{
		envoyutil.XFCCElement{By: "me"},
		envoyutil.XFCCElement{URI: "you"},
	}
	assert.Equal(expected, xfcc)

	// KV separator is the last character
	xfcc, err = envoyutil.ParseXFCC("by=cliffhanger;")
	assert.Equal(envoyutil.XFCC{}, xfcc)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	"Ends with separator character",
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	// Element separator is the last character
	xfcc, err = envoyutil.ParseXFCC("uri=cliffhanger,")
	assert.Equal(envoyutil.XFCC{}, xfcc)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:		envoyutil.ErrXFCCParsing,
		Message:	"Ends with separator character",
		StackTrace:	except.StackTrace,
	}
	assert.Equal(expectedErr, except)
	assert.Equal(envoyutil.XFCC{}, ele)

	// Empty header
	xfcc, err = envoyutil.ParseXFCC("")
	assert.Equal(envoyutil.XFCC{}, xfcc)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCC{}, ele)
}
