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
	fullXFCCTest                     = `By=spiffe://cluster.local/ns/blend/sa/yule;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688;Subject="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client";URI=spiffe://cluster.local/ns/blend/sa/cheer`
	xfccElementByTest                = `By=spiffe://cluster.local/ns/blend/sa/tide`
	xfccElementHashTest              = `Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688`
	xfccElementCertTest              = `Cert=` + xfccElementTestCertEncoded
	xfccElementChainTest             = `CHAIN=` + xfccElementTestCertEncoded
	xfccElementSubjectTest           = `SUBJECT="/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client"`
	xfccElementURITest               = `URI=spiffe://cluster.local/ns/blend/sa/quasar`
	xfccElementDNSTest               = `dns=http://frontend.lyft.com`
	xfccElementEndTest               = `dns=http://frontend.lyft.com;`
	xfccElementNoneTest              = `key=value;dns=http://frontend.lyft.com`
	xfccElementMultiTest             = `By=spiffe://cluster.local/ns/blend/sa/laser;Hash=468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688`
	xfccElementMalformedKeyTest      = `=value`
	xfccElementMultiMalformedKeyTest = `=value;dns=http://frontend.lyft.com`
	xfccElementMultiCertTest         = `cert=` + xfccElementTestCertEncoded + xfccElementTestCert

	xfccElementMalformedEncoding = "%"

	xfccElementTestCertEncoded = `-----BEGIN%20CERTIFICATE-----%0AMIIFKjCCAxICCQCA5%2FOCxg%2FqiDANBgkqhkiG9w0BAQsFADBXMQswCQYDVQQGEwJV%0AUzELMAkGA1UECAwCQ0ExFjAUBgNVBAcMDVNhbiBGcmFuY2lzY28xDTALBgNVBAoM%0ABEx5ZnQxFDASBgNVBAMMC1Rlc3QgQ2xpZW50MB4XDTIwMDYwNDE3NDkzNVoXDTIx%0AMDYwNDE3NDkzNVowVzELMAkGA1UEBhMCVVMxCzAJBgNVBAgMAkNBMRYwFAYDVQQH%0ADA1TYW4gRnJhbmNpc2NvMQ0wCwYDVQQKDARMeWZ0MRQwEgYDVQQDDAtUZXN0IENs%0AaWVudDCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBAKs6T8vcb8rIIkC4%0Aiz9h%2FOj6Iv%2BfazTLwNLK%2Fk58Ape5ZL0IdW6h8pWDlGnGz4X%2FTaJ5TwlamFo1h62v%0AsR8HPNOoLY0wmC2qHVquPF6eR9Lt5ejJiakr%2BYvf%2BU6LHXlOpOoot5rcTGoGBCf0%0AH3zmjdOE0o6hwJxMf54XQEVwNXqRrIDbY27mYS8eAVcSMrPUQVZ%2B3Vk1S56Imybz%0Adegi79IIoc6TzE5M7ChfJZBNNNZT08haJe6Oi%2FIgZhK3IexssY%2BQyD5uBSc7Mpas%0A6TstzeevIbeFy3Od2GhUy2Hz98qW%2FoO5iuerEArkNs4lB0J%2F0ARPHUDnmmH%2BqWYF%0APKealq2yEyXHHXrhDcSK%2FN5R64pp%2FVrxEas1qG20%2FCG4rixv36UJuEz5oUKNWyaR%0A268EI5Vecw%2BpK%2F0XC2%2Bhra9T%2FeP9JH0Fp43x7bdpQoxph8ZJZBsjbgCFMonf3ku1%0A9n74%2FxwvV6B0wp5C8jpwbGa85n%2BT8hogtO78mnpvxhTVJ7TOy596tI2apJ02edtD%0AJgsJV9MfZ%2FfGu3QZ6yN3rKVMPkZfC18cK04xy%2BroPo756CHkUHP5cz%2BKtJ7%2B8COR%0ArPDPxKBLOqwaSFcanQNONFIrffnZciiisCxjMHGoM4%2Fuix5gStlDC9%2FM5yyHt9He%0AldC8xL%2FyIalsa9Df7SL59Fd7T2JrAgMBAAEwDQYJKoZIhvcNAQELBQADggIBAGTb%0AOTddb6Yr37JigDYGKxSiidWfPgxTFr3DWTTE2aBhTi8K1Gr8Br1FAGRg4P114kxm%0AQBx3TxuCZcssMX2Z3GbL20DrHNExCl7hM%2FZVA2OCVhwfXGOHRQ6lcpmeQISWDNsN%0Atanlap%2FAgqKN%2F6At6JEYmuTSJnKc4Bfgk2GP5LPa63yJOlyvFb8ovKsCgb1ppVyw%0ARE%2B7AmB2DfDdVql4nHsDh5UBZRgVxMZ6xGnkYKaAUDKl4slejvKwXuzu2Xf%2BAd74%0AgjdLHzP0WmHlAggR5LIv%2F9xlvrsKCrNDDxWwOGeYk2WZl%2Fybud0RFKhLIqbbeMy7%0ADcdy04cJcqa9qRHYySgaWtM6Ab%2Fx9CJqdzR2NQZNnLgk6Vc3%2BoDjXMUuyM17WJAS%0ArenwJvanXvF9P1yPMByJQlXxkUehkCa%2FPs7E1O%2F%2BE2FJnvrtGVdYVR8Otbec1osS%0AmtJC6k7rgMhgvk63sCqQqaZwRWwLl2R5XcDZknUiqDKjuVHHA01II7jtGB1oyEIH%0Asp%2FrQlLNeyYlyhAlc3MhF5hu6nUjH%2B2%2BDuIHJsM0mEF0rjlbnp4bKJ%2FgF1COAIAL%0APzu2qAC%2BaOFldCmRonqUluayv6fQaQCeeh8sW2IjNVjA2ynKn2ybGIXH4mrH0KVa%0AJmUY%2B1YGMn7qbeHTma33N28Ec7hK%2BWByul746Nro%0A-----END%20CERTIFICATE-----`
	xfccElementTestCert        = `-----BEGIN CERTIFICATE-----
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
		By          string
		Expected    *url.URL
		ErrorRegexp *regexp.Regexp
	}
	testCases := []testCase{
		{By: "", Expected: &url.URL{}},
		// NOTE: Regex needed to support error format changes from go1.13 to go1.14
		{By: "\n", ErrorRegexp: regexp.MustCompile(`(?m)^parse ("\\n"|\n): net/url: invalid control character in URL$`)},
		{
			By: "spiffe://cluster.local/ns/blend/sa/yule",
			Expected: &url.URL{
				Scheme: "spiffe",
				Host:   "cluster.local",
				Path:   "/ns/blend/sa/yule",
			},
		},
	}
	for _, tc := range testCases {
		xe := envoyutil.XFCCElement{By: tc.By}
		uri, err := xe.DecodeBy()
		assert.Equal(tc.Expected, uri)
		if tc.ErrorRegexp != nil {
			assert.True(tc.ErrorRegexp.MatchString(err.Error()))
		} else {
			assert.Nil(err)
		}
	}
}

func TestXFCCElementDecodeHash(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		Hash     string
		Expected []byte
		Error    string
	}
	testCases := []testCase{
		{Hash: "", Expected: []byte("")},
		{Hash: "41434944", Expected: []byte("ACID")},
		{Hash: "x", Expected: nil, Error: "encoding/hex: invalid byte: U+0078 'x'"},
	}
	for _, tc := range testCases {
		xe := envoyutil.XFCCElement{Hash: tc.Hash}
		decoded, err := xe.DecodeHash()
		assert.Equal(tc.Expected, decoded)
		if tc.Error != "" {
			assert.Equal(tc.Error, err.Error())
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
		Cert   string
		Parsed *x509.Certificate
		Error  string
	}
	testCases := []testCase{
		{Cert: ""},
		{Cert: xfccElementTestCertEncoded, Parsed: parsedCert[0]},
		{Cert: "%", Error: `invalid URL escape "%"`},
		{
			Cert:  "-----BEGIN CERTIFICATE-----\nnope\n-----END CERTIFICATE-----\n",
			Error: "asn1: syntax error: truncated tag or length",
		},
		{
			Cert:  url.QueryEscape(xfccElementTestCert + "\n" + xfccElementTestCert),
			Error: "Error Parsing X-Forwarded-Client-Cert; Incorrect number of certificates; expected 1 got 2",
		},
		{
			Cert:  xfccElementMultiCertTest,
			Error: "Error Parsing X-Forwarded-Client-Cert; Incorrect number of certificates; expected 1 got 0",
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
		Chain  string
		Parsed []*x509.Certificate
		Error  string
	}
	testCases := []testCase{
		{Chain: ""},
		{
			Chain:  url.QueryEscape(xfccElementTestCert + "\n" + xfccElementTestCert),
			Parsed: parsedCerts,
		},
		{Chain: "%", Error: `invalid URL escape "%"`},
		{
			Chain: "-----BEGIN CERTIFICATE-----\nnope\n-----END CERTIFICATE-----\n",
			Error: "asn1: syntax error: truncated tag or length",
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
		URI         string
		Expected    *url.URL
		ErrorRegexp *regexp.Regexp
	}
	testCases := []testCase{
		{URI: "", Expected: &url.URL{}},
		// NOTE: Regex needed to support error format changes from go1.13 to go1.14
		{URI: "\r", ErrorRegexp: regexp.MustCompile(`(?m)^parse ("\\r"|\r): net/url: invalid control character in URL$`)},
		{
			URI: "spiffe://cluster.local/ns/first/sa/furst",
			Expected: &url.URL{
				Scheme: "spiffe",
				Host:   "cluster.local",
				Path:   "/ns/first/sa/furst",
			},
		},
	}
	for _, tc := range testCases {
		xe := envoyutil.XFCCElement{URI: tc.URI}
		uri, err := xe.DecodeURI()
		assert.Equal(tc.Expected, uri)
		if tc.ErrorRegexp != nil {
			assert.True(tc.ErrorRegexp.MatchString(err.Error()))
		} else {
			assert.Nil(err)
		}
	}
}

func TestXFCCElementString(t *testing.T) {
	assert := sdkAssert.New(t)

	type testCase struct {
		Element  *envoyutil.XFCCElement
		Expected string
		Override *envoyutil.XFCCElement
	}
	testCases := []testCase{
		{Element: &envoyutil.XFCCElement{}, Expected: ""},
		{Element: &envoyutil.XFCCElement{By: "hi"}, Expected: "By=hi"},
		{Element: &envoyutil.XFCCElement{Hash: "1389ab1"}, Expected: "Hash=1389ab1"},
		{Element: &envoyutil.XFCCElement{Cert: "anything-goes"}, Expected: "Cert=anything-goes"},
		{Element: &envoyutil.XFCCElement{Chain: "anything-goes"}, Expected: "Chain=anything-goes"},
		{
			Element:  &envoyutil.XFCCElement{Subject: "OU=Blent/CN=Test Client"},
			Expected: `Subject="OU=Blent/CN=Test Client"`,
			Override: &envoyutil.XFCCElement{Subject: `"OU=Blent/CN=Test Client"`},
		},
		{Element: &envoyutil.XFCCElement{URI: "bye"}, Expected: "URI=bye"},
		{
			Element:  &envoyutil.XFCCElement{DNS: []string{"web.invalid", "bye.invalid"}},
			Expected: "DNS=web.invalid;DNS=bye.invalid",
		},
	}
	for _, tc := range testCases {
		asString := tc.Element.String()
		assert.Equal(tc.Expected, asString)
		parsed, err := envoyutil.ParseXFCCElement(asString)
		assert.Nil(err)
		if tc.Override != nil {
			assert.Equal(tc.Override, &parsed)
		} else {
			assert.Equal(tc.Element, &parsed)
		}
	}

	// NOTE: For quoting and unquoting, `ParseXFCCElement()` is not fully an
	//       inverse to `XFCCElement.String()`.
	xe := envoyutil.XFCCElement{By: "a,b=10", URI: `c; "then" again`}
	assert.Equal("By=\"a,b=10\";URI=\"c; \\\"then\\\" again\"", xe.String())
}

func TestParseXFCCElement(t *testing.T) {
	assert := sdkAssert.New(t)

	ele, err := envoyutil.ParseXFCCElement(xfccElementByTest)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{By: "spiffe://cluster.local/ns/blend/sa/tide"}, ele)

	ele, err = envoyutil.ParseXFCCElement(xfccElementByTest + ";" + xfccElementByTest)
	assert.NotNil(err)
	except, ok := err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)
	assert.Equal(envoyutil.XFCCElement{By: "spiffe://cluster.local/ns/blend/sa/tide"}, ele)

	ele, err = envoyutil.ParseXFCCElement(xfccElementHashTest)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{Hash: "468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688"}, ele)

	ele, err = envoyutil.ParseXFCCElement(xfccElementHashTest + ";" + xfccElementHashTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr := &ex.Ex{
		Class:      envoyutil.ErrXFCCParsing,
		Message:    `Key already encountered "hash"`,
		StackTrace: except.StackTrace,
	}
	assert.Equal(expectedErr, except)

	ele, err = envoyutil.ParseXFCCElement(xfccElementCertTest)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{Cert: xfccElementTestCertEncoded}, ele)

	ele, err = envoyutil.ParseXFCCElement(xfccElementCertTest + ";" + xfccElementCertTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:      envoyutil.ErrXFCCParsing,
		Message:    `Key already encountered "cert"`,
		StackTrace: except.StackTrace,
	}
	assert.Equal(expectedErr, except)

	ele, err = envoyutil.ParseXFCCElement(xfccElementChainTest)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{Chain: xfccElementTestCertEncoded}, ele)

	ele, err = envoyutil.ParseXFCCElement(xfccElementChainTest + ";" + xfccElementChainTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:      envoyutil.ErrXFCCParsing,
		Message:    `Key already encountered "chain"`,
		StackTrace: except.StackTrace,
	}
	assert.Equal(expectedErr, except)

	ele, err = envoyutil.ParseXFCCElement(xfccElementSubjectTest)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{Subject: `"/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client"`}, ele)

	ele, err = envoyutil.ParseXFCCElement(xfccElementSubjectTest + ";" + xfccElementSubjectTest)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	expectedErr = &ex.Ex{
		Class:      envoyutil.ErrXFCCParsing,
		Message:    `Key already encountered "subject"`,
		StackTrace: except.StackTrace,
	}
	assert.Equal(expectedErr, except)

	ele, err = envoyutil.ParseXFCCElement(xfccElementURITest)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{URI: "spiffe://cluster.local/ns/blend/sa/quasar"}, ele)

	ele, err = envoyutil.ParseXFCCElement(xfccElementURITest + ";" + xfccElementURITest)
	assert.NotNil(err)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)
	assert.Equal(envoyutil.XFCCElement{URI: "spiffe://cluster.local/ns/blend/sa/quasar"}, ele)

	ele, err = envoyutil.ParseXFCCElement(xfccElementDNSTest)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{DNS: []string{"http://frontend.lyft.com"}}, ele)

	ele, err = envoyutil.ParseXFCCElement("dns=web.invalid;dns=blend.local.invalid")
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{DNS: []string{"web.invalid", "blend.local.invalid"}}, ele)

	_, err = envoyutil.ParseXFCCElement(xfccElementNoneTest)
	assert.NotNil(err)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)

	ele, err = envoyutil.ParseXFCCElement(xfccElementMultiTest)
	assert.Nil(err)
	expected := envoyutil.XFCCElement{
		By:   "spiffe://cluster.local/ns/blend/sa/laser",
		Hash: "468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688",
	}
	assert.Equal(expected, ele)

	_, err = envoyutil.ParseXFCCElement(xfccElementMalformedKeyTest)
	assert.NotNil(err)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)

	_, err = envoyutil.ParseXFCCElement(xfccElementMultiMalformedKeyTest)
	assert.NotNil(err)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)

	ele, err = envoyutil.ParseXFCCElement(xfccElementEndTest)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{DNS: []string{"http://frontend.lyft.com"}}, ele)

	ele, err = envoyutil.ParseXFCCElement("cert=" + xfccElementMalformedEncoding)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{Cert: "%"}, ele)

	ele, err = envoyutil.ParseXFCCElement("chain=" + xfccElementMalformedEncoding)
	assert.Nil(err)
	assert.Equal(envoyutil.XFCCElement{Chain: "%"}, ele)

	ele, err = envoyutil.ParseXFCCElement("=;")
	assert.NotNil(err)
	except, ok = err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)
	assert.Equal(envoyutil.XFCCElement{}, ele)

	// Test empty subject
	ele, err = envoyutil.ParseXFCCElement(`By=spiffe://cluster.local/ns/blend/sa/protocol;Hash=52114972613efb0820c5e32bfee0f0ee2a84859f7169da6c222300ef852a1129;Subject="";URI=spiffe://cluster.local/ns/blend/sa/world`)
	assert.Nil(err)
	expected = envoyutil.XFCCElement{
		By:      "spiffe://cluster.local/ns/blend/sa/protocol",
		Hash:    "52114972613efb0820c5e32bfee0f0ee2a84859f7169da6c222300ef852a1129",
		Subject: `""`,
		URI:     "spiffe://cluster.local/ns/blend/sa/world",
	}
	assert.Equal(expected, ele)
}

func TestParseXFCC(t *testing.T) {
	assert := sdkAssert.New(t)

	xfcc, err := envoyutil.ParseXFCC(fullXFCCTest + "," + fullXFCCTest)
	assert.Nil(err)
	expected := envoyutil.XFCC{
		envoyutil.XFCCElement{
			By:      "spiffe://cluster.local/ns/blend/sa/yule",
			Hash:    "468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688",
			Subject: "/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client",
			URI:     "spiffe://cluster.local/ns/blend/sa/cheer",
		},
		envoyutil.XFCCElement{
			By:      "spiffe://cluster.local/ns/blend/sa/yule",
			Hash:    "468ed33be74eee6556d90c0149c1309e9ba61d6425303443c0748a02dd8de688",
			Subject: "/C=US/ST=CA/L=San Francisco/OU=Lyft/CN=Test Client",
			URI:     "spiffe://cluster.local/ns/blend/sa/cheer",
		},
	}
	assert.Equal(expected, xfcc)

	_, err = envoyutil.ParseXFCC(xfccElementMalformedKeyTest)
	assert.NotNil(err)
	except, ok := err.(*ex.Ex)
	assert.True(ok)
	assert.NotNil(except)
	assert.Equal(envoyutil.ErrXFCCParsing, except.Class)
}
