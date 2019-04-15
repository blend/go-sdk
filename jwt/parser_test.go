package jwt_test

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/jwt"
	"github.com/blend/go-sdk/jwt/test"
)

var keyFuncError error = fmt.Errorf("error loading key")

var (
	jwtTestDefaultKey *rsa.PublicKey
	defaultKeyFunc    jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return jwtTestDefaultKey, nil }
	emptyKeyFunc      jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return nil, nil }
	errorKeyFunc      jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return nil, keyFuncError }
	nilKeyFunc        jwt.Keyfunc = nil
)

func init() {
	jwtTestDefaultKey = test.MustLoadRSAPublicKey(test.SampleKeyPublic)
}

var jwtTestData = []struct {
	name        string
	tokenString string
	keyfunc     jwt.Keyfunc
	claims      jwt.Claims
	valid       bool
	err         error
	inner       error
	parser      *jwt.Parser
}{
	{
		"basic",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		true,
		nil,
		nil,
		nil,
	},
	{
		"basic expired",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "exp": float64(time.Now().Unix() - 100)},
		false,
		jwt.ErrValidation,
		jwt.ErrValidationExpired,
		nil,
	},
	{
		"basic nbf",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100)},
		false,
		jwt.ErrValidation,
		jwt.ErrValidationNotBefore,
		nil,
	},
	{
		"expired and nbf",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": float64(time.Now().Unix() + 100), "exp": float64(time.Now().Unix() - 100)},
		false,
		jwt.ErrValidation,
		jwt.ErrValidationExpired,
		nil,
	},
	{
		"basic invalid",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.EhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		jwt.ErrValidation,
		jwt.ErrValidationSignature,
		nil,
	},
	{
		"basic nokeyfunc",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		nilKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		jwt.ErrValidation,
		jwt.ErrKeyfuncUnset,
		nil,
	},
	{
		"basic nokey",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		emptyKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		jwt.ErrValidation,
		jwt.ErrValidationSignature,
		nil,
	},
	{
		"basic errorkey",
		"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIifQ.FhkiHkoESI_cG3NPigFrxEk9Z60_oXrOT2vGm9Pn6RDgYNovYORQmmA0zs1AoAOf09ly2Nx2YAg6ABqAYga1AcMFkJljwxTT5fYphTuqpWdy4BELeSYJx5Ty2gmr8e7RonuUztrdD5WfPqLKMm1Ozp_T6zALpRmwTIW0QPnaBXaQD90FplAg46Iy1UlDKr-Eupy0i5SLch5Q-p2ZpaL_5fnTIUDlxC3pWhJTyx_71qDI-mAA_5lE_VdroOeflG56sSmDxopPEG3bFlSu1eowyBfxtu0_CuVd-M42RU75Zc4Gsj6uV77MBtbMrf4_7M_NUTSgoIF3fRqxrj0NzihIBg",
		errorKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		keyFuncError,
		nil,
		nil,
	},
	{
		"invalid signing method",
		"",
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		false,
		jwt.ErrValidation,
		jwt.ErrInvalidSigningMethod,
		&jwt.Parser{ValidMethods: []string{"HS256"}},
	},
	{
		"valid signing method",
		"",
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar"},
		true,
		nil,
		nil,
		&jwt.Parser{ValidMethods: []string{"RS256", "HS256"}},
	},
	{
		"JSON Number",
		"",
		defaultKeyFunc,
		jwt.MapClaims{"foo": json.Number("123.4")},
		true,
		nil,
		nil,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"Standard Claims",
		"",
		defaultKeyFunc,
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 10).Unix(),
		},
		true,
		nil,
		nil,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"JSON Number - basic expired",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "exp": json.Number(fmt.Sprintf("%v", time.Now().Unix()-100))},
		false,
		jwt.ErrValidation,
		jwt.ErrValidationExpired,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"JSON Number - basic nbf",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100))},
		false,
		jwt.ErrValidation,
		jwt.ErrValidationNotBefore,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"JSON Number - expired and nbf",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100)), "exp": json.Number(fmt.Sprintf("%v", time.Now().Unix()-100))},
		false,
		jwt.ErrValidation,
		jwt.ErrValidationExpired,
		&jwt.Parser{UseJSONNumber: true},
	},
	{
		"SkipClaimsValidation during token parsing",
		"", // autogen
		defaultKeyFunc,
		jwt.MapClaims{"foo": "bar", "nbf": json.Number(fmt.Sprintf("%v", time.Now().Unix()+100))},
		true,
		nil,
		nil,
		&jwt.Parser{UseJSONNumber: true, SkipClaimsValidation: true},
	},
}

func TestParser_Parse(t *testing.T) {
	privateKey := test.MustLoadRSAPrivateKey(test.SampleKey)

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		// If the token string is blank, use helper function to generate string
		if data.tokenString == "" {
			data.tokenString = test.MakeSampleToken(data.claims, privateKey)
		}

		// Parse the token
		var token *jwt.Token
		var err error
		var parser = data.parser
		if parser == nil {
			parser = new(jwt.Parser)
		}
		// Figure out correct claims type
		switch data.claims.(type) {
		case jwt.MapClaims:
			token, err = parser.ParseWithClaims(data.tokenString, jwt.MapClaims{}, data.keyfunc)
		case *jwt.StandardClaims:
			token, err = parser.ParseWithClaims(data.tokenString, &jwt.StandardClaims{}, data.keyfunc)
		}

		// Verify result matches expectation
		if !reflect.DeepEqual(data.claims, token.Claims) {
			t.Errorf("[%v] Claims mismatch. Expecting: %v Got: %v", data.name, data.claims, token.Claims)
		}

		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
		}

		if !data.valid && err == nil {
			t.Errorf("[%v] Invalid token passed validation", data.name)
		}

		if (err == nil && !token.Valid) || (err != nil && token.Valid) {
			t.Errorf("[%v] Inconsistent behavior between returned error and token.Valid", data.name)
		}

		if data.err != nil {
			if err == nil {
				t.Errorf("[%v] Expecting error. Didn't get one.", data.name)
			} else {
				if !ex.Is(err, data.err) {
					t.Errorf("[%v] Errors don't match expectation. `%v` != `%v`", data.name, err, data.err)
				}
				if data.inner != nil {
					if typed, ok := err.(*ex.Ex); ok {
						if !ex.Is(typed.Inner, data.inner) {
							t.Errorf("[%v] Causing Errors don't match expectation. `%+v` != `%+v`", data.name, typed.Inner, data.inner)
						}
					} else {
						t.Errorf("[%v] Errors aren't exceptions and we expect a causing error", data.name)
					}
				}
			}
		}
		if data.valid && token.Signature == "" {
			t.Errorf("[%v] Signature is left unpopulated after parsing", data.name)
		}
	}
}

func TestParser_ParseUnverified(t *testing.T) {
	privateKey := test.MustLoadRSAPrivateKey(test.SampleKey)

	// Iterate over test data set and run tests
	for _, data := range jwtTestData {
		// If the token string is blank, use helper function to generate string
		if data.tokenString == "" {
			data.tokenString = test.MakeSampleToken(data.claims, privateKey)
		}

		// Parse the token
		var token *jwt.Token
		var err error
		var parser = data.parser
		if parser == nil {
			parser = new(jwt.Parser)
		}
		// Figure out correct claims type
		switch data.claims.(type) {
		case jwt.MapClaims:
			token, _, err = parser.ParseUnverified(data.tokenString, jwt.MapClaims{})
		case *jwt.StandardClaims:
			token, _, err = parser.ParseUnverified(data.tokenString, &jwt.StandardClaims{})
		}

		if err != nil {
			t.Errorf("[%v] Invalid token", data.name)
		}

		// Verify result matches expectation
		if !reflect.DeepEqual(data.claims, token.Claims) {
			t.Errorf("[%v] Claims mismatch. Expecting: %v  Got: %v", data.name, data.claims, token.Claims)
		}

		if data.valid && err != nil {
			t.Errorf("[%v] Error while verifying token: %T:%v", data.name, err, err)
		}
	}
}

// Helper method for benchmarking various methods
func benchmarkSigning(b *testing.B, method jwt.SigningMethod, key interface{}) {
	t := jwt.New(method)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := t.SignedString(key); err != nil {
				b.Fatal(err)
			}
		}
	})

}
