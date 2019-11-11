package oauth

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/util"
	"github.com/blend/go-sdk/webutil"
)

func TestNew(t *testing.T) {
	assert := assert.New(t)
	assert.Empty(New().Secret())
}

func TestNewFromConfig(t *testing.T) {
	assert := assert.New(t)

	m, err := NewFromConfig(&Config{
		RedirectURI:  "https://app.com/oauth/google",
		HostedDomain: "foo.com",
		ClientID:     "foo_client",
		ClientSecret: "bar_secret",
	})

	assert.Nil(err)
	assert.Empty(m.Secret())
	assert.Equal("https://app.com/oauth/google", m.RedirectURI())
	assert.Equal("foo_client", m.ClientID())
	assert.Equal("bar_secret", m.ClientSecret())
}

func TestNewFromConfigWithSecret(t *testing.T) {
	assert := assert.New(t)

	m, err := NewFromConfig(&Config{
		Secret: base64.StdEncoding.EncodeToString([]byte("test string")),
	})

	assert.Nil(err)
	assert.NotEmpty(m.Secret())
	assert.Equal("test string", string(m.Secret()))
}

func TestManagerOAuthURLWithFullyQualifiedRedirectURI(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithClientID("test_client_id").
		WithHostedDomain("test.blend.com").
		WithRedirectURI("https://local.shortcut-service.centrio.com/oauth/google")

	oauthURL, err := m.OAuthURL(nil)
	assert.Nil(err)

	parsed, err := url.Parse(oauthURL)
	assert.Nil(err)
	assert.Equal("test.blend.com", parsed.Query().Get("hd"), "we should set the hosted domain if it's configured")
}

func TestManagerOAuthURL(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithClientID("test_client_id").
		WithRedirectURI("/oauth/google")

	oauthURL, err := m.OAuthURL(&http.Request{RequestURI: "https://test.blend.com/foo"})
	assert.Nil(err)

	_, err = url.Parse(oauthURL)
	assert.Nil(err)
}

func TestManagerGetRedirectURI(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithClientID("test_client_id").
		WithRedirectURI("/oauth/google")

	redirectURI := m.getRedirectURI(&http.Request{Proto: "https", Host: "test.blend.com", Header: http.Header{webutil.HeaderXForwardedProto: {webutil.SchemeHTTPS}}})
	parsedRedirectURI, err := url.Parse(redirectURI)
	assert.Nil(err)
	assert.Equal("https", parsedRedirectURI.Scheme)
	assert.Equal("test.blend.com", parsedRedirectURI.Host)
	assert.Equal("/oauth/google", parsedRedirectURI.Path)
}

func TestManagerGetRedirectURIFullyQualified(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithClientID("test_client_id").
		WithRedirectURI("https://test.blend.com/oauth/google")

	redirectURI := m.getRedirectURI(nil)

	parsedRedirectURI, err := url.Parse(redirectURI)
	assert.Nil(err)
	assert.Equal("https", parsedRedirectURI.Scheme)
	assert.Equal("test.blend.com", parsedRedirectURI.Host)
	assert.Equal("/oauth/google", parsedRedirectURI.Path)
}

func TestManagerGetRedirectURIFullyQualifiedHTTP(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithClientID("test_client_id").
		WithRedirectURI("http://test.blend.com/oauth/google")

	redirectURI := m.getRedirectURI(nil)

	parsedRedirectURI, err := url.Parse(redirectURI)
	assert.Nil(err)
	assert.Equal("http", parsedRedirectURI.Scheme)
	assert.Equal("test.blend.com", parsedRedirectURI.Host)
	assert.Equal("/oauth/google", parsedRedirectURI.Path)
}

func TestManagerGetRedirectURIFullyQualifiedSPDY(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithClientID("test_client_id").
		WithRedirectURI("spdy://test.blend.com/oauth/google")

	redirectURI := m.getRedirectURI(nil)
	parsedRedirectURI, err := url.Parse(redirectURI)
	assert.Nil(err)
	assert.Equal("spdy", parsedRedirectURI.Scheme)
	assert.Equal("test.blend.com", parsedRedirectURI.Host)
	assert.Equal("/oauth/google", parsedRedirectURI.Path)
}

func TestManagerOAuthURLRedirect(t *testing.T) {
	assert := assert.New(t)

	m := New().
		WithClientID("test_client_id").
		WithRedirectURI("https://local.shortcut-service.centrio.com/oauth/google")

	urlFragment, err := m.OAuthURL(nil, "bar_foo")
	assert.Nil(err)

	u, err := url.Parse(urlFragment)
	assert.Nil(err)
	assert.NotEmpty(u.Query().Get("state"))

	state := u.Query().Get("state")
	deserialized, err := DeserializeState(state)
	assert.Nil(err)
	assert.Nil(m.ValidateState(deserialized))
	assert.Equal("bar_foo", deserialized.RedirectURL)
}

func TestManagerValidateProfile(t *testing.T) {
	assert := assert.New(t)

	blender := &Profile{
		Email: "bailey@blend.com",
	}

	personal := &Profile{
		Email: "bailey@gmail.com",
	}

	suffixMatch := &Profile{
		Email: "bailey@sailblend.com",
	}

	prefixMatch := &Profile{
		Email: "bailey@blend.com.au",
	}

	empty := New()
	assert.Nil(empty.ValidateProfile(blender), "we should not error if the hosted domain is not configured")

	hosted := New().WithHostedDomain("blend.com")
	assert.Nil(hosted.ValidateProfile(blender), "we should pass for @blend.com")
	assert.NotNil(hosted.ValidateProfile(personal), "we fail for non-@blend.com emails")
	assert.NotNil(hosted.ValidateProfile(suffixMatch), "we fail for non-@blend.com emails")
	assert.NotNil(hosted.ValidateProfile(prefixMatch), "we fail for non-@blend.com emails")

	hostedPrefixed := New().WithHostedDomain("@blend.com")
	assert.Nil(hostedPrefixed.ValidateProfile(blender), "we should pass for @blend.com")
	assert.NotNil(hostedPrefixed.ValidateProfile(personal), "we fail for non-@blend.com emails")
	assert.NotNil(hostedPrefixed.ValidateProfile(suffixMatch), "we fail for non-@blend.com emails")
	assert.NotNil(hostedPrefixed.ValidateProfile(prefixMatch), "we fail for non-@blend.com emails")
}

func TestManagerValidateState(t *testing.T) {
	assert := assert.New(t)

	insecure := New()
	assert.Nil(insecure.ValidateState(insecure.CreateState()))

	secure := New().WithSecret(util.Crypto.MustCreateKey(32))
	assert.Nil(secure.ValidateState(secure.CreateState()))
}
