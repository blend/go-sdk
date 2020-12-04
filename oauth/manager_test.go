package oauth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/crypto"
	"github.com/blend/go-sdk/jwt"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/webutil"
)

func Test_Manager_Finish(t *testing.T) {
	it := assert.New(t)

	pk0, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk0pem))
	it.Nil(err)
	pk1, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk1pem))
	it.Nil(err)
	keys := []jwt.JWK{
		createJWK(pk0),
		createJWK(pk1),
	}
	keysResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.Header().Set(http.CanonicalHeaderKey("Cache-Control"), "public, max-age=23196, must-revalidate, no-transform") // set cache control
		rw.Header().Set(http.CanonicalHeaderKey("Expires"), time.Now().UTC().AddDate(0, 1, 0).Format(http.TimeFormat))    // set expires
		rw.Header().Set(http.CanonicalHeaderKey("Date"), time.Now().UTC().Format(http.TimeFormat))                        // set date
		rw.WriteHeader(200)
		json.NewEncoder(rw).Encode(struct {
			Keys []jwt.JWK `json:"keys"`
		}{
			Keys: keys,
		})
	}))
	defer keysResponder.Close()

	codeResponse, err := createCodeResponse("test_client_id", keys[1].KID, pk1)
	it.Nil(err)

	codeResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.WriteHeader(200)
		rw.Write(codeResponse)
	}))
	defer codeResponder.Close()

	profileResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if accessToken := req.Header.Get(webutil.HeaderAuthorization); accessToken != "Bearer test_access_token" {
			http.Error(rw, "not authorized", http.StatusUnauthorized)
			return
		}

		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.WriteHeader(200)
		fmt.Fprintf(rw, `{
			"id": "12012312390931",
			"email": "bailey@test.blend.com",
			"verified_email": true,
			"name": "Bailey Dog",
			"given_name": "Bailey",
			"family_name": "Dog",
			"picture": "https://example.com/bailey.jpg",
			"locale": "en",
			"hd": "test.blend.com"
		  }`)
	}))
	defer profileResponder.Close()

	mgr, err := New(
		OptClientID("test_client_id"),
		OptClientSecret(crypto.MustCreateKeyString(32)),
		OptSecret(crypto.MustCreateKey(32)),
		OptAllowedDomains("test.blend.com"),
	)
	it.Nil(err)
	mgr.PublicKeyCache.FetchPublicKeysDefaults = []r2.Option{
		r2.OptURL(keysResponder.URL),
	}
	mgr.FetchProfileDefaults = []r2.Option{
		r2.OptURL(profileResponder.URL),
	}
	mgr.Endpoint = oauth2.Endpoint{
		AuthStyle: oauth2.AuthStyleInParams,
		TokenURL:  codeResponder.URL,
	}
	finishRequest := &http.Request{
		URL: &url.URL{
			RawQuery: (url.Values{
				"code":  []string{"test_code"},
				"state": []string{MustSerializeState(mgr.CreateState())},
			}).Encode(),
		},
	}

	res, err := mgr.Finish(finishRequest)
	it.Nil(err)
	it.Equal("bailey@test.blend.com", res.Profile.Email)
	it.Equal("Bailey", res.Profile.GivenName)
	it.Equal("Dog", res.Profile.FamilyName)
	it.Equal("en", res.Profile.Locale)
	it.Equal("https://example.com/bailey.jpg", res.Profile.PictureURL)
}

func Test_Manager_Finish_DisallowedDomain(t *testing.T) {
	it := assert.New(t)

	pk0, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk0pem))
	it.Nil(err)
	pk1, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk1pem))
	it.Nil(err)
	keys := []jwt.JWK{
		createJWK(pk0),
		createJWK(pk1),
	}
	keysResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.Header().Set(http.CanonicalHeaderKey("Cache-Control"), "public, max-age=23196, must-revalidate, no-transform") // set cache control
		rw.Header().Set(http.CanonicalHeaderKey("Expires"), time.Now().UTC().AddDate(0, 1, 0).Format(http.TimeFormat))    // set expires
		rw.Header().Set(http.CanonicalHeaderKey("Date"), time.Now().UTC().Format(http.TimeFormat))                        // set date
		rw.WriteHeader(200)
		json.NewEncoder(rw).Encode(struct {
			Keys []jwt.JWK `json:"keys"`
		}{
			Keys: keys,
		})
	}))
	defer keysResponder.Close()

	codeResponse, err := createCodeResponse("test_client_id", keys[1].KID, pk1)
	it.Nil(err)

	codeResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.WriteHeader(200)
		rw.Write(codeResponse)
	}))
	defer codeResponder.Close()

	profileResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if accessToken := req.URL.Query().Get("access_token"); accessToken != "test_access_token" {
			http.Error(rw, "not authorized", http.StatusUnauthorized)
			return
		}

		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.WriteHeader(200)
		fmt.Fprintf(rw, `{
			"id": "12012312390931",
			"email": "bailey@test.blend.com",
			"verified_email": true,
			"name": "Bailey Dog",
			"given_name": "Bailey",
			"family_name": "Dog",
			"picture": "https://example.com/bailey.jpg",
			"locale": "en",
			"hd": "test.blend.com"
		  }`)
	}))
	defer profileResponder.Close()

	mgr, err := New(
		OptClientID("test_client_id"),
		OptClientSecret(crypto.MustCreateKeyString(32)),
		OptSecret(crypto.MustCreateKey(32)),
		OptAllowedDomains("blend.com"),
	)
	it.Nil(err)
	mgr.PublicKeyCache.FetchPublicKeysDefaults = []r2.Option{
		r2.OptURL(keysResponder.URL),
	}
	mgr.FetchProfileDefaults = []r2.Option{
		r2.OptURL(profileResponder.URL),
	}
	mgr.Endpoint = oauth2.Endpoint{
		AuthStyle: oauth2.AuthStyleInParams,
		TokenURL:  codeResponder.URL,
	}
	finishRequest := &http.Request{
		URL: &url.URL{
			RawQuery: (url.Values{
				"code":  []string{"test_code"},
				"state": []string{MustSerializeState(mgr.CreateState())},
			}).Encode(),
		},
	}

	res, err := mgr.Finish(finishRequest)
	it.NotNil(err)
	it.Empty(res.Profile.Email)
}

func Test_Manager_Finish_FailsAudience(t *testing.T) {
	it := assert.New(t)

	pk0, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk0pem))
	it.Nil(err)
	pk1, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk1pem))
	it.Nil(err)
	keys := []jwt.JWK{
		createJWK(pk0),
		createJWK(pk1),
	}
	keysResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.Header().Set(http.CanonicalHeaderKey("Cache-Control"), "public, max-age=23196, must-revalidate, no-transform") // set cache control
		rw.Header().Set(http.CanonicalHeaderKey("Expires"), time.Now().UTC().AddDate(0, 1, 0).Format(http.TimeFormat))    // set expires
		rw.Header().Set(http.CanonicalHeaderKey("Date"), time.Now().UTC().Format(http.TimeFormat))                        // set date
		rw.WriteHeader(200)
		json.NewEncoder(rw).Encode(struct {
			Keys []jwt.JWK `json:"keys"`
		}{
			Keys: keys,
		})
	}))
	defer keysResponder.Close()

	codeResponse, err := createCodeResponse("not_test_client_id", keys[1].KID, pk1)
	it.Nil(err)

	codeResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.WriteHeader(200)
		rw.Write(codeResponse)
	}))
	defer codeResponder.Close()

	mgr, err := New(
		OptClientID("test_client_id"),
		OptClientSecret(crypto.MustCreateKeyString(32)),
		OptSecret(crypto.MustCreateKey(32)),
		OptAllowedDomains("blend.com"),
	)
	it.Nil(err)
	mgr.PublicKeyCache.FetchPublicKeysDefaults = []r2.Option{
		r2.OptURL(keysResponder.URL),
	}
	mgr.Endpoint = oauth2.Endpoint{
		AuthStyle: oauth2.AuthStyleInParams,
		TokenURL:  codeResponder.URL,
	}
	finishRequest := &http.Request{
		URL: &url.URL{
			RawQuery: (url.Values{
				"code":  []string{"test_code"},
				"state": []string{MustSerializeState(mgr.CreateState())},
			}).Encode(),
		},
	}

	res, err := mgr.Finish(finishRequest)
	it.NotNil(err)
	it.Empty(res.Profile.Email)
}

func Test_Manager_Finish_FailsVerification(t *testing.T) {
	it := assert.New(t)

	pk0, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk0pem))
	it.Nil(err)
	pk1, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk1pem))
	it.Nil(err)
	pk2, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pk2pem))
	it.Nil(err)
	keys := []jwt.JWK{
		createJWK(pk0),
		createJWK(pk1),
	}
	keysResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.Header().Set(http.CanonicalHeaderKey("Cache-Control"), "public, max-age=23196, must-revalidate, no-transform") // set cache control
		rw.Header().Set(http.CanonicalHeaderKey("Expires"), time.Now().UTC().AddDate(0, 1, 0).Format(http.TimeFormat))    // set expires
		rw.Header().Set(http.CanonicalHeaderKey("Date"), time.Now().UTC().Format(http.TimeFormat))                        // set date
		rw.WriteHeader(200)
		json.NewEncoder(rw).Encode(struct {
			Keys []jwt.JWK `json:"keys"`
		}{
			Keys: keys,
		})
	}))
	defer keysResponder.Close()

	codeResponse, err := createCodeResponse("test_client_id", uuid.V4().String(), pk2)
	it.Nil(err)

	codeResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.WriteHeader(200)
		rw.Write(codeResponse)
	}))
	defer codeResponder.Close()

	profileResponder := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if accessToken := req.URL.Query().Get("access_token"); accessToken != "test_access_token" {
			http.Error(rw, "not authorized", http.StatusUnauthorized)
			return
		}

		rw.Header().Set(http.CanonicalHeaderKey("Content-Type"), "application/json; charset=UTF-8")
		rw.WriteHeader(200)
		fmt.Fprintf(rw, `{
			"id": "12012312390931",
			"email": "bailey@test.blend.com",
			"verified_email": true,
			"name": "Bailey Dog",
			"given_name": "Bailey",
			"family_name": "Dog",
			"picture": "https://example.com/bailey.jpg",
			"locale": "en",
			"hd": "test.blend.com"
		  }`)
	}))
	defer profileResponder.Close()

	mgr, err := New(
		OptClientID("test_client_id"),
		OptClientSecret(crypto.MustCreateKeyString(32)),
		OptSecret(crypto.MustCreateKey(32)),
		OptAllowedDomains("test.blend.com"),
	)
	it.Nil(err)
	mgr.PublicKeyCache.FetchPublicKeysDefaults = []r2.Option{
		r2.OptURL(keysResponder.URL),
	}
	mgr.FetchProfileDefaults = []r2.Option{
		r2.OptURL(profileResponder.URL),
	}
	mgr.Endpoint = oauth2.Endpoint{
		AuthStyle: oauth2.AuthStyleInParams,
		TokenURL:  codeResponder.URL,
	}
	finishRequest := &http.Request{
		URL: &url.URL{
			RawQuery: (url.Values{
				"code":  []string{"test_code"},
				"state": []string{MustSerializeState(mgr.CreateState())},
			}).Encode(),
		},
	}

	res, err := mgr.Finish(finishRequest)
	it.NotNil(err)
	it.Empty(res.Profile.Email)
}

func Test_MustNew(t *testing.T) {
	assert := assert.New(t)
	assert.Empty(MustNew().Secret)
	assert.NotEmpty(MustNew().Endpoint.AuthURL)
	assert.NotEmpty(MustNew().Scopes)
}

func Test_NewFromConfig(t *testing.T) {
	assert := assert.New(t)

	m, err := New(OptConfig(Config{
		RedirectURI:  "https://app.com/oauth/google",
		HostedDomain: "foo.com",
		ClientID:     "foo_client",
		ClientSecret: "bar_secret",
	}))

	assert.Nil(err)
	assert.Empty(m.Secret)
	assert.Equal("https://app.com/oauth/google", m.RedirectURL)
	assert.Equal("foo_client", m.ClientID)
	assert.Equal("bar_secret", m.ClientSecret)
}

func Test_NewFromConfigWithSecret(t *testing.T) {
	assert := assert.New(t)

	m, err := New(OptConfig(Config{
		Secret: base64.StdEncoding.EncodeToString([]byte("test string")),
	}))

	assert.Nil(err)
	assert.NotEmpty(m.Secret)
	assert.Equal("test string", string(m.Secret))
}

func Test_Manager_OAuthURL_FullyQualifiedRedirectURI(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.HostedDomain = "test.blend.com"
	m.RedirectURL = "https://local.shortcut-service.centrio.com/oauth/google"

	oauthURL, err := m.OAuthURL(nil)
	assert.Nil(err)

	parsed, err := url.Parse(oauthURL)
	assert.Nil(err)
	assert.Equal("test.blend.com", parsed.Query().Get("hd"), "we should set the hosted domain if it's configured")
}

func Test_Manager_OAuthURL(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.RedirectURL = "/oauth/google"

	oauthURL, err := m.OAuthURL(&http.Request{RequestURI: "https://test.blend.com/foo"})
	assert.Nil(err)

	_, err = url.Parse(oauthURL)
	assert.Nil(err)
}

func Test_Manager_OAuthURLRedirect(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.RedirectURL = "https://local.shortcut-service.centrio.com/oauth/google"

	urlFragment, err := m.OAuthURL(nil, OptStateRedirectURI("bar_foo"))
	assert.Nil(err)

	u, err := url.Parse(urlFragment)
	assert.Nil(err)
	assert.NotEmpty(u.Query().Get("state"))

	state := u.Query().Get("state")
	deserialized, err := DeserializeState(state)
	assert.Nil(err)
	assert.Nil(m.ValidateState(deserialized))
	assert.Equal("bar_foo", deserialized.RedirectURI)
}

func Test_Manager_ValidateState(t *testing.T) {
	assert := assert.New(t)

	insecure := MustNew()
	assert.Nil(insecure.ValidateState(insecure.CreateState()))

	secure := MustNew()
	secure.Secret = crypto.MustCreateKey(32)
	assert.Nil(secure.ValidateState(secure.CreateState()))

	wrongKey := MustNew()
	wrongKey.Secret = crypto.MustCreateKey(32)

	assert.NotNil(secure.ValidateState(wrongKey.CreateState()))
}
