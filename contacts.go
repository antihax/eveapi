package eveapi

import "fmt"

type Contact struct {
	Standing    float64 `json:"standing,omitempty"`
	ContactType string  `json:"contactType,omitempty"`
	Contact     struct {
		Href string `json:"href,omitempty"`
		Name string `json:"name,omitempty"`
		ID   int64  `json:"id,omitempty"`
	} `json:"contact,omitempty"`
	Watched bool `json:"watched,omitempty"`
}

func (c *AuthenticatedClient) SetContact(id int64, standing float64) (*VerifyResponse, error) {
	v, err := c.Verify()
	if err != nil {
		return nil, err
	}

	contact := Contact{Standing: standing}
	contact.Contact.ID = id

	url := fmt.Sprintf(c.base.CREST+"characters/%d/contacts/", c.characterID)
	ret := make(map[string]interface{})
	res, err := c.doJSON("POST", url, contact, &ret)
	fmt.Printf("%s %+v %+v %v\n", url, ret, res, err)
	if err != nil {
		return nil, err
	}
	return v, nil
}
