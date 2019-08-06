package oauth

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/blend/go-sdk/crypto"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/webutil"
)

func TestMustNew(t *testing.T) {
	assert := assert.New(t)
	assert.Empty(MustNew().Secret)
}

func TestNewFromConfig(t *testing.T) {
	assert := assert.New(t)

	m, err := New(OptConfig(Config{
		RedirectURI:  "https://app.com/oauth/google",
		HostedDomain: "foo.com",
		ClientID:     "foo_client",
		ClientSecret: "bar_secret",
	}))

	assert.Nil(err)
	assert.Empty(m.Secret)
	assert.Equal("https://app.com/oauth/google", m.RedirectURI)
	assert.Equal("foo_client", m.ClientID)
	assert.Equal("bar_secret", m.ClientSecret)
}

func TestNewFromConfigWithSecret(t *testing.T) {
	assert := assert.New(t)

	m, err := New(OptConfig(Config{
		Secret: base64.StdEncoding.EncodeToString([]byte("test string")),
	}))

	assert.Nil(err)
	assert.NotEmpty(m.Secret)
	assert.Equal("test string", string(m.Secret))
}

func TestManagerOAuthURLWithFullyQualifiedRedirectURI(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.HostedDomain = "test.blend.com"
	m.RedirectURI = "https://local.shortcut-service.centrio.com/oauth/google"

	oauthURL, err := m.OAuthURL(nil)
	assert.Nil(err)

	parsed, err := url.Parse(oauthURL)
	assert.Nil(err)
	assert.Equal("test.blend.com", parsed.Query().Get("hd"), "we should set the hosted domain if it's configured")
}

func TestManagerOAuthURL(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.RedirectURI = "/oauth/google"

	oauthURL, err := m.OAuthURL(&http.Request{RequestURI: "https://test.blend.com/foo"})
	assert.Nil(err)

	_, err = url.Parse(oauthURL)
	assert.Nil(err)
}

func TestManagerGetRedirectURI(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.RedirectURI = "/oauth/google"

	redirectURI := m.getRedirectURI(&http.Request{Proto: "spdy", Host: "test.blend.com", Header: http.Header{webutil.HeaderXForwardedProto: {webutil.SchemeHTTPS}}})
	parsedRedirectURI, err := url.Parse(redirectURI)
	assert.Nil(err)
	assert.Equal(webutil.SchemeHTTPS, parsedRedirectURI.Scheme)
	assert.Equal("test.blend.com", parsedRedirectURI.Host)
	assert.Equal("/oauth/google", parsedRedirectURI.Path)
}

func TestManagerGetRedirectURIFullyQualified(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.RedirectURI = "https://test.blend.com/oauth/google"

	redirectURI := m.getRedirectURI(nil)

	parsedRedirectURI, err := url.Parse(redirectURI)
	assert.Nil(err)
	assert.Equal("https", parsedRedirectURI.Scheme)
	assert.Equal("test.blend.com", parsedRedirectURI.Host)
	assert.Equal("/oauth/google", parsedRedirectURI.Path)
}

func TestManagerGetRedirectURIFullyQualifiedHTTP(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.RedirectURI = "http://test.blend.com/oauth/google"

	redirectURI := m.getRedirectURI(nil)

	parsedRedirectURI, err := url.Parse(redirectURI)
	assert.Nil(err)
	assert.Equal("http", parsedRedirectURI.Scheme)
	assert.Equal("test.blend.com", parsedRedirectURI.Host)
	assert.Equal("/oauth/google", parsedRedirectURI.Path)
}

func TestManagerGetRedirectURIFullyQualifiedSPDY(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.RedirectURI = "spdy://test.blend.com/oauth/google"

	redirectURI := m.getRedirectURI(nil)
	parsedRedirectURI, err := url.Parse(redirectURI)
	assert.Nil(err)
	assert.Equal("spdy", parsedRedirectURI.Scheme)
	assert.Equal("test.blend.com", parsedRedirectURI.Host)
	assert.Equal("/oauth/google", parsedRedirectURI.Path)
}

func TestManagerOAuthURLRedirect(t *testing.T) {
	assert := assert.New(t)

	m, err := New()
	assert.Nil(err)
	m.ClientID = "test_client_id"
	m.RedirectURI = "https://local.shortcut-service.centrio.com/oauth/google"

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

	empty := MustNew()
	assert.Nil(empty.ValidateProfile(blender), "we should not error if the hosted domain is not configured")

	hosted := MustNew()
	hosted.HostedDomain = "blend.com"
	assert.Nil(hosted.ValidateProfile(blender), "we should pass for @blend.com")
	assert.NotNil(hosted.ValidateProfile(personal), "we fail for non-@blend.com emails")
	assert.NotNil(hosted.ValidateProfile(suffixMatch), "we fail for non-@blend.com emails")
	assert.NotNil(hosted.ValidateProfile(prefixMatch), "we fail for non-@blend.com emails")

	hostedPrefixed := MustNew()
	hostedPrefixed.HostedDomain = "@blend.com"
	assert.Nil(hostedPrefixed.ValidateProfile(blender), "we should pass for @blend.com")
	assert.NotNil(hostedPrefixed.ValidateProfile(personal), "we fail for non-@blend.com emails")
	assert.NotNil(hostedPrefixed.ValidateProfile(suffixMatch), "we fail for non-@blend.com emails")
	assert.NotNil(hostedPrefixed.ValidateProfile(prefixMatch), "we fail for non-@blend.com emails")
}

func TestManagerValidateState(t *testing.T) {
	assert := assert.New(t)

	insecure := MustNew()
	assert.Nil(insecure.ValidateState(insecure.CreateState()))

	secure := MustNew()
	secure.Secret = crypto.MustCreateKey(32)
	assert.Nil(secure.ValidateState(secure.CreateState()))
}

func TestManagerRequestDefaulkts(t *testing.T) {
	assert := assert.New(t)

	mgr := MustNew(
		OptRequestDefaults(
			r2.OptHeaderValue("foo", "bar"),
		),
	)
	assert.NotEmpty(mgr.RequestDefaults)

}
