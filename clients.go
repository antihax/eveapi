package eveapi

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

// AnonymousClient for Public CREST and Public XML API queries.
type AnonymousClient struct {
	httpClient *http.Client
	base       EveURI
	userAgent  string
}

// AuthenticatedClient for Private CREST and Private XML API queries. SSO Authenticated.
type AuthenticatedClient struct {
	AnonymousClient
	token     *oauth2.Token
	character *VerifyResponse
}

// ErrorMessage format if a CREST query fails.
type ErrorMessage struct {
	Message string
}

// Executes a request generated with newRequest
func (c *AnonymousClient) executeRequest(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Creates a new http.Request for a public resource.
func (c *AnonymousClient) newRequest(method, urlStr string, body interface{}, mediaType string) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, rel.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", BASE_API_VERSION)
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

// Provides a new http.Request for an authenticated resource.
func (c *AuthenticatedClient) newRequest(method, urlStr string, body interface{}, mediaType string) (*http.Request, error) {
	req, err := c.AnonymousClient.newRequest(method, urlStr, body, mediaType)
	if err != nil {
		return nil, err
	}

	c.token.SetAuthHeader(req)
	return req, nil
}

// Calls a resource from the public XML API
func (c *AnonymousClient) doXML(method, urlStr string, body interface{}, v interface{}) (*http.Response, error) {
	req, err := c.newRequest(method, urlStr, body, "application/xml")
	if err != nil {
		return nil, err
	}
	<-xmlThrottle // Throttle XML requests
	res, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err := xml.Unmarshal([]byte(buf), v); err != nil {
		return nil, err
	}
	return res, nil
}

// Calls a resource from the public CREST
func (c *AnonymousClient) doJSON(method, urlStr string, body interface{}, v interface{}, mediaType string) (*http.Response, error) {
	req, err := c.newRequest(method, urlStr, body, mediaType)
	if err != nil {
		return nil, err
	}
	<-anonThrottle              // Throttle Anonymous CREST requests
	anonConnectionLimit <- true // Limit concurrent requests
	res, err := c.executeRequest(req)
	<-anonConnectionLimit

	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		e := &ErrorMessage{}
		if err := json.Unmarshal([]byte(buf), e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	}
	if err := json.Unmarshal([]byte(buf), v); err != nil {
		return nil, err
	}

	return res, nil
}

// Calls a resource from authenticated CREST.
func (c *AuthenticatedClient) doJSON(method, urlStr string, body interface{}, v interface{}, mediaType string) (*http.Response, error) {
	req, err := c.newRequest(method, urlStr, body, mediaType)
	if err != nil {
		return nil, err
	}

	<-authedThrottle // Throttle Authenticated CREST requests

	res, err := c.executeRequest(req)
	if err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusCreated {
		return res, nil
	} else if res.StatusCode != http.StatusOK {
		e := &ErrorMessage{}
		if err = json.Unmarshal([]byte(buf), e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	} else {
		if err = json.Unmarshal([]byte(buf), v); err != nil {
			return res, err
		}
	}
	return res, err
}

// SetUI set the user agent string of the CREST and XML client.
// It is recommended to change this so that CCP can identify your app.
func (c *AnonymousClient) SetUA(userAgent string) {
	c.userAgent = userAgent
}

// UseCustomURL allows the base URLs to be changed should the need arise
// for a third party proxy to be used.
func (c *AnonymousClient) UseCustomURL(custom EveURI) {
	c.base = custom
}

// GetCharacterID returns the associated characterID for this authenticated client.
// Verify must be called prior to this becoming available or it will return 0.
func (c *AuthenticatedClient) GetCharacterID() int64 {
	if c.character == nil {
		return 0
	}
	return c.character.CharacterID
}

// UseTestServer forces this client to use the test server URLs.
func (c *AnonymousClient) UseTestServer(testServer bool) {
	if testServer == true {
		c.base = eveSisi
	} else {
		c.base = eveTQ
	}
}

// NewAnonymousClient generates a new anonymous client.
// Caller must provide a caching http.Client that obeys all cacheUntil timers
// One Anonymous client per IP address or rate limits will be exceeded resulting in a ban.
func NewAnonymousClient(client *http.Client) *AnonymousClient {
	c := &AnonymousClient{}
	c.base = eveTQ
	c.httpClient = client
	c.userAgent = USER_AGENT
	return c
}

type VerifyResponse struct {
	CharacterID        int64
	CharacterName      string
	ExpiresOn          string
	Scopes             string
	TokenType          string
	CharacterOwnerHash string
}

// Verify the client and collect user information.
func (c *AuthenticatedClient) Verify() (*VerifyResponse, error) {
	v := &VerifyResponse{}
	_, err := c.doJSON("GET", c.base.Login+"oauth/verify", nil, v, "application/json;")
	c.character = v
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Verify that the client is validated.
func (c *AuthenticatedClient) validateClient() error {
	var err error
	if c.character != nil {
		return nil
	}
	c.character, err = c.Verify()
	if err != nil {
		return err
	}
	return nil
}
