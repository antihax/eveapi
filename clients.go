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

type AnonymousClient struct {
	httpClient *http.Client
	base       EveURI
	userAgent  string
}

type AuthenticatedClient struct {
	AnonymousClient
	token     *oauth2.Token
	character *VerifyResponse
}
type ErrorMessage struct {
	Message string
}

const (
	userAgent = "https://github.com/antihax/eveapi"
	mediaType = "application/json"
)

func (c *AnonymousClient) executeRequest(req *http.Request) (*http.Response, error) {
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

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

func (c *AuthenticatedClient) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	req, err := c.AnonymousClient.newRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}

	c.token.SetAuthHeader(req)
	return req, nil
}

func (c *AnonymousClient) doXML(method, urlStr string, body interface{}, v interface{}) (*http.Response, error) {
	req, err := c.newRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
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

func (c *AnonymousClient) doJSON(method, urlStr string, body interface{}, v interface{}) (*http.Response, error) {
	req, err := c.newRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	res, err := c.executeRequest(req)
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

func (c *AuthenticatedClient) doXML(method, urlStr string, body interface{}, v interface{}) (*http.Response, error) {
	req, err := c.newRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
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

func (c *AuthenticatedClient) doJSON(method, urlStr string, body interface{}, v interface{}) (*http.Response, error) {
	req, err := c.newRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
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
		if err := json.Unmarshal([]byte(buf), e); err != nil {
			return nil, err
		}
		return nil, errors.New(e.Message)
	} else {
		if err := json.Unmarshal([]byte(buf), v); err != nil {
			return res, err
		}
	}
	return res, err
}

func (c *AnonymousClient) SetUA(userAgent string) {
	c.userAgent = userAgent
}

func (c *AnonymousClient) UseCustomURL(custom EveURI) {
	c.base = custom
}

func (c *AuthenticatedClient) GetCharacterID() int64 {
	return c.character.CharacterID
}

func (c *AnonymousClient) UseTestServer(testServer bool) {
	if testServer == true {
		c.base = eveSisi
	} else {
		c.base = eveTQ
	}
}

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

// Verify the client and collect user information.
func (c *AuthenticatedClient) Verify() (*VerifyResponse, error) {
	v := &VerifyResponse{}
	_, err := c.doJSON("GET", c.base.Login+"oauth/verify", nil, v)
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
