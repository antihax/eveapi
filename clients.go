package eveapi

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"
)

type AnonymousClient struct {
	httpClient *http.Client
	base       EveURI
}

type AuthenticatedClient struct {
	AnonymousClient
	tokenSource oauth2.TokenSource
}

func (c *AnonymousClient) UseCustomURL(custom EveURI) {
	c.base = custom
}

func (c *AnonymousClient) UseTestServer(testServer bool) {
	if testServer == true {
		c.base = eveSisi
	} else {
		c.base = eveTQ
	}
}

// NewAuthenticatedClient assigns a token to a client.
func NewAnonymousClient(client *http.Client) *AnonymousClient {
	c := &AnonymousClient{}
	c.base = eveTQ
	c.httpClient = client
	return c
}

// NewAuthenticatedClient assigns a token to a client.
func NewAuthenticatedClient(client *http.Client, tok CRESTTokenP) *AuthenticatedClient {
	c := &AuthenticatedClient{}
	c.base = eveTQ
	c.tokenSource = oauth2.StaticTokenSource(tok)
	c.httpClient = oauth2.NewClient(createContext(client), c.tokenSource)
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

func decodeJSON(res *http.Response, ret interface{}) error {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(buf), ret); err != nil {
		return err
	}
	return err
}

func decodeXML(res *http.Response, ret interface{}) error {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := xml.Unmarshal([]byte(buf), ret); err != nil {
		return err
	}
	return err
}

func (c *AuthenticatedClient) Verify() (*VerifyResponse, error) {
	v := &VerifyResponse{}
	res, err := c.httpClient.Get(c.base.Login + "/oauth/verify")
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	err = decodeJSON(res, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
