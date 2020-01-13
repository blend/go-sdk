package oauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"github.com/blend/go-sdk/ex"
	"github.com/blend/go-sdk/r2"
	"github.com/blend/go-sdk/stringutil"
	"github.com/blend/go-sdk/uuid"
	"github.com/blend/go-sdk/webutil"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// New returns a new manager mutated by a given set of options.
func New(options ...Option) (*Manager, error) {
	manager := &Manager{}
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
	FetchProfileDefaults []r2.Option
	Tracer               Tracer
	Secret               []byte
	Scopes               []string
	RedirectURI          string
	HostedDomain         string
	ClientID             string
	ClientSecret         string
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
	oauthURL = m.conf(r).AuthCodeURL(state, opts...)
	return
}

// Finish processes the returned code, exchanging for an access token, and fetches the user profile.
func (m *Manager) Finish(r *http.Request) (result *Result, err error) {
	if m.Tracer != nil {
		tf := m.Tracer.Start(r.Context(), m.conf(r))
		if tf != nil {
			defer func() { tf.Finish(r.Context(), m.conf(r), result, err) }()
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
	tok, err := m.conf(r).Exchange(r.Context(), code)
	if err != nil {
		err = ex.New(ErrFailedCodeExchange, ex.OptInner(err))
		return
	}

	result.Response.AccessToken = tok.AccessToken
	result.Response.TokenType = tok.TokenType
	result.Response.RefreshToken = tok.RefreshToken
	result.Response.Expiry = tok.Expiry

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
	res, err := r2.New("https://www.googleapis.com/oauth2/v1/userinfo", append(m.FetchProfileDefaults,
		r2.OptGet(),
		r2.OptContext(ctx),
		r2.OptQueryValue("alt", "json"),
		r2.OptQueryValue("access_token", accessToken),
	)...).Do()

	if err != nil {
		return
	}
	if res.StatusCode > 299 {
		err = ex.New(ErrGoogleResponseStatus, ex.OptMessagef("status code: %d", res.StatusCode))
		return
	}
	defer res.Body.Close()
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

// ValidateProfile validates a profile.
func (m *Manager) ValidateProfile(p *Profile) error {
	if len(m.HostedDomain) == 0 {
		return nil
	}

	workingDomain := m.HostedDomain
	if !strings.HasPrefix(workingDomain, "@") {
		workingDomain = fmt.Sprintf("@%s", workingDomain)
	}
	if !stringutil.HasSuffixCaseless(p.Email, workingDomain) {
		return ErrInvalidHostedDomain
	}
	return nil
}

// --------------------------------------------------------------------------------
// internal helpers
// --------------------------------------------------------------------------------

func (m *Manager) conf(r *http.Request) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     m.ClientID,
		ClientSecret: m.ClientSecret,
		RedirectURL:  m.getRedirectURI(r),
		Scopes:       m.Scopes,
		Endpoint:     google.Endpoint,
	}
}

func (m *Manager) getRedirectURI(r *http.Request) string {
	if stringutil.HasPrefixCaseless(m.RedirectURI, "https://") ||
		stringutil.HasPrefixCaseless(m.RedirectURI, "http://") ||
		stringutil.HasPrefixCaseless(m.RedirectURI, "spdy://") {
		return m.RedirectURI
	}
	requestURI := &url.URL{
		Scheme: webutil.GetProto(r),
		Host:   webutil.GetHost(r),
		Path:   m.RedirectURI,
	}
	return requestURI.String()
}

func (m *Manager) hash(plaintext string) string {
	return base64.URLEncoding.EncodeToString(m.hmac([]byte(plaintext)))
}

// hmac hashes data with the given key.
func (m *Manager) hmac(plainText []byte) []byte {
	mac := hmac.New(sha512.New, m.Secret)
	mac.Write([]byte(plainText))
	return mac.Sum(nil)
}
