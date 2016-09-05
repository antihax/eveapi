/*
Package eveapi impliments a CREST and XML api rate limited client for EVE's
public and private APIs.

Caching

Caching is not implimented by the client and thus it is required to utilize
a caching http client. It is highly recommended to utilize a client capable
of caching the entire cluster of API clients.

An example using gregjones/httpcache and memcache:

	import (
		"github.com/bradfitz/gomemcache/memcache"
		"github.com/gregjones/httpcache"
		httpmemcache "github.com/gregjones/httpcache/memcache"
	)

	func main() {
		// Connect to the memcache server
		cache := memcache.New(MemcachedAddresses...)

		// Create a memcached http client for the CCP APIs.
		transport := httpcache.NewTransport(httpmemcache.NewWithClient(cache))
		transport.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment}
		client = &http.Client{Transport: transport}

		// Get our EVE API.
		eve := eveapi.NewAnonymousClient(client)
	}

Anonymous Client and Public Endpoints

All public endpoints are available through a simple anonymous client. It
is recommended to use this for all public calls due to the higher rate limits.

Anonymous clients are created simply by supplying a caching HTTP Client.

	eve := eveapi.NewAnonymousClient(client)

Private Clients

Private clients are generated through an SSO Authenticator and Tokens (both Access and Refresh).
They require the application clientID, secretKey, and redirect URL. These must match exactly
to what was provided to CCP on the Manage Applications page.

	scopes := []string{eveapi.ScopeCharacterContactsRead,
		eveapi.ScopeCharacterContactsWrite}

	tokenAuthenticator = eveapi.NewSSOAuthenticator(clientID,
		secretKey,redirectURL,scopes)

	privateClient := tokenAuthenticator.GetClientFromToken(client, token)

One authenticator can spawn as many clients as needed at once, each with their own tokens.

SSO

Obtaining tokens for client requires two HTTP handlers. One to generate and redirect
to the SSO URL, and one to receive the response. It is manditory to create a random
state and compare this state on return to prevent token injection attacks on the application.

In the example custom httpHandlers below, SSOAuthenticator, is a created by NewSSOAuthenticator, with scopes.

	func eveSSO(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
		// Generate a random state string
		b := make([]byte, 16)
		rand.Read(b)
		state := base64.URLEncoding.EncodeToString(b)

		// Save the state on the session
		s.Values["state"] = state
		err := s.Save(r, w)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// Generate the SSO URL with the state string
		url := c.SSOAuthenticator.AuthorizeURL(state, true)

		// Send the user to the URL
		http.Redirect(w, r, url, 302)
		return http.StatusMovedPermanently, nil
	}

	func eveSSOAnswer(c *appContext.AppContext, w http.ResponseWriter, r *http.Request, s *sessions.Session) (int, error) {
		// get our code and state
		code := r.FormValue("code")
		state := r.FormValue("state")

		// Verify the state matches our randomly generated string from earlier.
		if s.Values["state"] != state {
			return http.StatusInternalServerError, errors.New("Invalid State.")
		}

		// Exchange the code for an Access and Refresh token.
		token, err := c.SSOAuthenticator.TokenExchange(c.HTTPClient, code)
		if err != nil {
			return http.StatusInternalServerError, errors.New("Failed Token Exchange")
		}

		// Obtain a private authenticated client for this token.
		cli := c.SSOAuthenticator.GetClientFromToken(c.HTTPClient, token)

		// Verify the access token (will automaticlly use Refresh token if it is expired)
		v, err := cli.Verify()
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// Save the verification structure on the session for quick access.
		s.Values["character"] = v
		err = s.Save(r, w)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		// Redirect to the account page.
		http.Redirect(w, r, "/account", 302)
		return http.StatusMovedPermanently, nil
	}

*/
package eveapi
