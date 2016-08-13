package eveapi

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type AnonymousClient struct {
	httpClient *http.Client
	base       EveURI
	userAgent  string
}

type AuthenticatedClient struct {
	AnonymousClient
	tokenSource oauth2.TokenSource
	characterID int64
}

const (
	userAgent = "https://github.com/antihax/eveapi"
	mediaType = "application/json"
)

func (c *AnonymousClient) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, rel.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", "application/vnd.ccp.eve.Api-v3+json")
	req.Header.Add("User-Agent", c.userAgent)
	return req, nil
}

func (c *AnonymousClient) do(method, urlStr string, body interface{}) (*http.Response, error) {
	r, err := c.newRequest(method, urlStr, body)
	res, err := c.httpClient.Do(r)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *AnonymousClient) doXML(method, urlStr string, body interface{}, v interface{}) (*http.Response, error) {
	res, err := c.do(method, urlStr, body)
	defer res.Body.Close()
	err = decodeXML(res, v)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *AnonymousClient) doJSON(method, urlStr string, body interface{}, v interface{}) (*http.Response, error) {
	res, err := c.do(method, urlStr, body)
	defer res.Body.Close()
	err = decodeJSON(res, v)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *AnonymousClient) SetUA(userAgent string) {
	c.userAgent = userAgent
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
	c.userAgent = userAgent
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
	_, err := c.doJSON("GET", c.base.Login+"/oauth/verify", nil, v)
	c.characterID = v.CharacterID
	if err != nil {
		return nil, err
	}
	return v, nil
}
