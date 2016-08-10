package eveapi

func (c *AuthenticatedClient) SetContact() (*VerifyResponse, error) {
	v := &VerifyResponse{}
	_, err := c.doJSON("GET", c.base.Login+"/oauth/verify", nil, v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
