package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2/clientcredentials"

	"golang.org/x/oauth2"
)

type GetContactsPayload struct {
	ID string `json:"id"`
}

type Email struct {
	FromName  string
	FromEmail string
	ToName    string
	ToEmail   string
	Subject   string
	Message   string
}

type GoogleService struct {
	ctx context.Context
}

// reuseTokenSource is a TokenSource that holds a single token in memory
// and validates its expiry before each call to retrieve it with
// Token. If it's expired, it will be auto-refreshed using the
// new TokenSource.
type ReuseTokenSource struct {
	new oauth2.TokenSource // called when t is expired.

	mu sync.Mutex // guards t
	t  *oauth2.Token
}

type GoogleConfig struct {
	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// TokenURL is the resource server's token endpoint
	// URL. This is a constant specific to each server.
	TokenURL string

	// Scope specifies optional requested permissions.
	Scopes []string

	OAuth2Config      *oauth2.Config
	ClientCredsConfig *clientcredentials.Config
}

// Token returns the current token if it's still valid, else will
// refresh the current token (using r.Context for HTTP client
// information) and return the new one.
func (s *ReuseTokenSource) Token() (*oauth2.Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	t, err := s.new.Token()
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}

// tokenRefresher is a TokenSource that makes "grant_type"=="refresh_token"
// HTTP requests to renew a token using a RefreshToken.
type TokenRefresher struct {
	ctx          context.Context // used to get HTTP requests
	conf         *GoogleConfig
	refreshToken string
}

// WARNING: Token is not safe for concurrent access, as it
// updates the tokenRefresher's refreshToken field.
// Within this package, it is used by reuseTokenSource which
// synchronizes calls to this method with its own mutex.
func (tf *TokenRefresher) Token() (*oauth2.Token, error) {
	log.Printf("Entered Token() for TokenRefresher with refresh token: %v", tf.refreshToken)
	tk, err := retrieveToken(tf.ctx, tf.conf, url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {tf.refreshToken},
		"client_id":     {tf.conf.ClientID},
		"client_secret": {tf.conf.ClientSecret},
	}, tf.refreshToken)

	if err != nil {
		return nil, err
	}
	return tk, err
}

type tokenRespBody struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // in seconds
}

// retrieveToken takes a *Config and uses that to retrieve an *internal.Token.
// This token is then mapped from *internal.Token into an *oauth2.Token which is returned along
// with an error..
func retrieveToken(ctx context.Context, c *GoogleConfig, v url.Values, refreshToken string) (*oauth2.Token, error) {
	req, _ := http.NewRequest("POST", c.TokenURL, strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.ClientID, c.ClientSecret)

	var r *http.Response
	var err error
	if c.OAuth2Config != nil {
		r, err = c.OAuth2Config.Client(ctx, &oauth2.Token{RefreshToken: refreshToken}).Do(req)
	} else {
		r, err = c.ClientCredsConfig.Client(ctx).Do(req)
	}

	if err != nil {
		log.Printf("retrieve token triggered error: %v", err.Error())
		return nil, err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	if c := r.StatusCode; c < 200 || c > 299 {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v\nResponse: %s", r.Status, body)
	}
	resp := &tokenRespBody{}
	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch content {
	case "application/x-www-form-urlencoded", "text/plain":
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		resp.AccessToken = vals.Get("access_token")
		resp.TokenType = vals.Get("token_type")
		resp.RefreshToken = vals.Get("refresh_token")
		resp.ExpiresIn, _ = strconv.ParseInt(vals.Get("expires_in"), 10, 64)
	default:
		if err = json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}
	}
	token := &oauth2.Token{
		AccessToken:  resp.AccessToken,
		TokenType:    resp.TokenType,
		RefreshToken: resp.RefreshToken,
	}
	// Don't overwrite `RefreshToken` with an empty value
	// if this was a token refreshing request.
	if resp.RefreshToken == "" {
		token.RefreshToken = v.Get("refresh_token")
	}
	if resp.ExpiresIn == 0 {
		token.Expiry = time.Time{}
	} else {
		token.Expiry = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	}
	return token, nil
}

// providerAuthHeaderWorks reports whether the OAuth2 server identified by the tokenURL
// implements the OAuth2 spec correctly
// See https://code.google.com/p/goauth2/issues/detail?id=31 for background.
// In summary:
// - Reddit only accepts client secret in the Authorization header
// - Dropbox accepts either it in URL param or Auth header, but not both.
// - Google only accepts URL param (not spec compliant?), not Auth header
func providerAuthHeaderWorks(tokenURL string) bool {
	if strings.HasPrefix(tokenURL, "https://accounts.google.com/") ||
		strings.HasPrefix(tokenURL, "https://github.com/") ||
		strings.HasPrefix(tokenURL, "https://api.instagram.com/") ||
		strings.HasPrefix(tokenURL, "https://www.douban.com/") ||
		strings.HasPrefix(tokenURL, "https://api.dropbox.com/") ||
		strings.HasPrefix(tokenURL, "https://api.soundcloud.com/") ||
		strings.HasPrefix(tokenURL, "https://www.linkedin.com/") {
		// Some sites fail to implement the OAuth2 spec fully.
		return false
	}
	// Assume the provider implements the spec properly
	// otherwise. We can add more exceptions as they're
	// discovered. We will _not_ be adding configurable hooks
	// to this package to let users select server bugs.
	return true
}

func NewGoogleService() *GoogleService {
	return &GoogleService{}
}
