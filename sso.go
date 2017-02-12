package eveapi

import (
	"encoding/json"
	"net/http"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type contextOAuth2Key *oauth2.Token

var ContextOAuth2 contextOAuth2Key

// SSOAuthenticator provides interfacing to the CREST SSO. NewSSOAuthenticator is used to create
// this structure.

// [TODO] lose this mutex and allow scopes to change without conflict.
type SSOAuthenticator struct {
	httpClient *http.Client
	// Hide this...
	oauthConfig *oauth2.Config
	scopeLock   sync.Mutex
}

// Redirect type to hide oauth2 API
type CRESTToken oauth2.Token

type CRESTTokenSource oauth2.TokenSource

// NewSSOAuthenticator create a new CREST SSO Authenticator.
// Requires your application clientID, clientSecret, and redirectURL.
// RedirectURL must match exactly to what you registered with CCP.
func NewSSOAuthenticator(client *http.Client, clientID string, clientSecret string, redirectURL string, scopes []string) *SSOAuthenticator {

	if client == nil {
		return nil
	}

	c := &SSOAuthenticator{}

	c.httpClient = client

	c.oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://login.eveonline.com/oauth/authorize",
			TokenURL: "https://login.eveonline.com/oauth/token",
		},
		Scopes:      scopes,
		RedirectURL: redirectURL,
	}
	return c
}

// AuthorizeURL returns a url for an end user to authenticate with EVE SSO
// and return success to the redirectURL.
// It is important to create a significatly unique state for this request
// and verify the state matches when returned to the redirectURL.
func (c *SSOAuthenticator) AuthorizeURL(state string, onlineAccess bool, scopes []string) string {
	var url string

	// lock so we cannot use another requests scopes by racing.
	c.scopeLock.Lock()

	// Save the default scopes.
	saveScopes := c.oauthConfig.Scopes
	if scopes != nil {
		c.oauthConfig.Scopes = scopes
	}

	// Generate the URL
	if onlineAccess == true {
		url = c.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOnline)
	} else {
		url = c.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	}

	// Return the scopes
	c.oauthConfig.Scopes = saveScopes

	// Unlock mutex. [TODO] This is seriously hacky... need to fix
	c.scopeLock.Unlock()

	return url
}

// TokenExchange exchanges the code returned to the redirectURL with
// the CREST server to an access token. A caching client must be passed.
// This client MUST cache per CCP guidelines or face banning.
func (c *SSOAuthenticator) TokenExchange(code string) (*CRESTToken, error) {

	tok, err := c.oauthConfig.Exchange(createContext(c.httpClient), code)
	if err != nil {
		return nil, err
	}
	return (*CRESTToken)(tok), nil
}

// TokenSource creates a refreshable token that can be passed to ESI functions
func (c *SSOAuthenticator) TokenSource(token *CRESTToken) (CRESTTokenSource, error) {
	return (CRESTTokenSource)(c.oauthConfig.TokenSource(createContext(c.httpClient), (*oauth2.Token)(token))), nil
}

// Add custom clients to the context.
func createContext(httpClient *http.Client) context.Context {
	parent := oauth2.NoContext
	ctx := context.WithValue(parent, oauth2.HTTPClient, httpClient)
	return ctx
}

// TokenToJSON helper function to convert a token to a storable format.
func TokenToJSON(token *CRESTToken) (string, error) {
	if d, err := json.Marshal(token); err != nil {
		return "", err
	} else {
		return string(d), nil
	}
}

// TokenFromJSON helper function to convert stored JSON to a token.
func TokenFromJSON(jsonStr string) (*CRESTToken, error) {
	var token CRESTToken
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
