package oauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/webutil"
)

// New returns a new manager mutated by a given set of options.
func New(options ...Option) (*Manager, error) {
	manager := &Manager{
		Config: oauth2.Config{
			Endpoint: google.Endpoint,
			Scopes:   DefaultScopes,
		},
		PublicKeyCache: new(PublicKeyCache),
	}

	for _, option := range options {
		if err := option(manager); err != nil {
			return nil, err
		}
	}
	return manager, nil
}

// MustNew returns a new manager mutated by a given set of options
// and will panic on error.
func MustNew(options ...Option) *Manager {
	m, err := New(options...)
	if err != nil {
		panic(err)
	}
	return m
}

// Manager is the oauth manager.
type Manager struct {
	oauth2.Config
	Tracer Tracer

	Secret []byte

	HostedDomain   string
	AllowedDomains []string

	FetchProfileDefaults []r2.Option
	PublicKeyCache       *PublicKeyCache
}

// OAuthURL is the auth url for google with a given clientID.
// This is typically the link that a user will click on to start the auth process.
func (m *Manager) OAuthURL(r *http.Request, stateOptions ...StateOption) (oauthURL string, err error) {
	var state string
	state, err = SerializeState(m.CreateState(stateOptions...))
	if err != nil {
		return
	}
	var opts []oauth2.AuthCodeOption
	if len(m.HostedDomain) > 0 {
		opts = append(opts, oauth2.SetAuthURLParam("hd", m.HostedDomain))
	}
	oauthURL = m.AuthCodeURL(state, opts...)
	return
}

// Finish processes the returned code, exchanging for an access token, and fetches the user profile.
func (m *Manager) Finish(r *http.Request) (result *Result, err error) {
	if m.Tracer != nil {
		tf := m.Tracer.Start(r.Context(), &m.Config)
		if tf != nil {
			defer func() { tf.Finish(r.Context(), &m.Config, result, err) }()
		}
	}

	// grab the code off the request.
	code := r.URL.Query().Get("code")
	if len(code) == 0 {
		err = ErrCodeMissing
		return
	}

	// fetch the state
	state := r.URL.Query().Get("state")
	result = &Result{}
	if len(state) > 0 {
		var deserialized State
		deserialized, err = DeserializeState(state)
		if err != nil {
			return
		}
		result.State = deserialized
	}
	err = m.ValidateState(result.State)
	if err != nil {
		return
	}

	// Handle the exchange code to initiate a transport.
	tok, err := m.Exchange(r.Context(), code)
	if err != nil {
		err = ex.New(ErrFailedCodeExchange, ex.OptInner(err))
		return
	}

	jwtClaims, err := ParseTokenJWT(tok, m.PublicKeyCache.Keyfunc(r.Context()))
	if err != nil {
		err = ex.New(ErrInvalidJWT, ex.OptInner(err))
		return
	}
	if err = m.ValidateJWT(jwtClaims); err != nil {
		return
	}

	result.Response.AccessToken = tok.AccessToken
	result.Response.TokenType = tok.TokenType
	result.Response.RefreshToken = tok.RefreshToken
	result.Response.Expiry = tok.Expiry
	result.Response.HostedDomain = jwtClaims.HD

	var prof Profile
	prof, err = m.FetchProfile(r.Context(), tok.AccessToken)
	if err != nil {
		return
	}
	result.Profile = prof
	return
}

// FetchProfile gets a google profile for an access token.
func (m *Manager) FetchProfile(ctx context.Context, accessToken string) (profile Profile, err error) {
	res, err := r2.New("https://www.googleapis.com/oauth2/v1/userinfo", append([]r2.Option{
		r2.OptGet(),
		r2.OptContext(ctx),
		r2.OptQueryValue("alt", "json"),
		r2.OptHeaderValue(webutil.HeaderAuthorization, fmt.Sprintf("Bearer %s", accessToken)),
	}, m.FetchProfileDefaults...)...).Do()
	if err != nil {
		return
	}
	defer res.Body.Close()
	if code := res.StatusCode; code < 200 || code > 299 {
		err = ex.New(ErrGoogleResponseStatus, ex.OptMessagef("status code: %d", res.StatusCode))
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&profile); err != nil {
		err = ex.New(ErrProfileJSONUnmarshal, ex.OptInner(err))
		return
	}
	return
}

// CreateState creates auth state.
func (m *Manager) CreateState(options ...StateOption) (state State) {
	for _, opt := range options {
		opt(&state)
	}
	if len(m.Secret) > 0 && state.Token == "" && state.SecureToken == "" {
		state.Token = uuid.V4().String()
		state.SecureToken = m.hash(state.Token)
	}
	return
}

// --------------------------------------------------------------------------------
// Validation Helpers
// --------------------------------------------------------------------------------

// ValidateState validates oauth state.
func (m *Manager) ValidateState(state State) error {
	if len(m.Secret) > 0 {
		expected := m.hash(state.Token)
		actual := state.SecureToken
		if !hmac.Equal([]byte(expected), []byte(actual)) {
			return ErrInvalidAntiforgeryToken
		}
	}
	return nil
}

// ValidateJWT returns if the jwt is valid or not.
func (m *Manager) ValidateJWT(jwtClaims *GoogleClaims) error {
	if jwtClaims.Audience != m.ClientID {
		return ex.New(ErrInvalidJWTAudience, ex.OptMessagef("audience: %s", jwtClaims.Audience))
	}
	if jwtClaims.Issuer != GoogleIssuer && jwtClaims.Issuer != GoogleIssuerAlternate {
		return ex.New(ErrInvalidJWTIssuer, ex.OptMessagef("issuer: %s", jwtClaims.Issuer))
	}
	if len(m.AllowedDomains) > 0 {
		if strings.TrimSpace(jwtClaims.HD) == "" {
			return ex.New(ErrInvalidJWTHostedDomain, ex.OptMessagef("hosted domain: likely gmail.com, but empty"))
		}
		var matchedDomain bool
		for _, domain := range m.AllowedDomains {
			if strings.EqualFold(domain, jwtClaims.HD) {
				matchedDomain = true
				break
			}
		}
		if !matchedDomain {
			return ex.New(ErrInvalidJWTHostedDomain, ex.OptMessagef("hosted domain: %s", jwtClaims.HD))
		}
	}
	return nil
}

// --------------------------------------------------------------------------------
// internal helpers
// --------------------------------------------------------------------------------

func (m *Manager) hash(plaintext string) string {
	return base64.URLEncoding.EncodeToString(m.hmac([]byte(plaintext)))
}

// hmac hashes data with the given key.
func (m *Manager) hmac(plainText []byte) []byte {
	mac := hmac.New(sha512.New, m.Secret)
	_, _ = mac.Write([]byte(plainText))
	return mac.Sum(nil)
}
