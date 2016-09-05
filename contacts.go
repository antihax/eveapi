package eveapi

import "fmt"

const contactCreateType = "application/vnd.ccp.eve.ContactCreate-v1"

type ContactCreateV1 struct {
	Standing    float64 `json:"standing,omitempty"`
	ContactType string  `json:"contactType,omitempty"`
	Contact     struct {
		Href string `json:"href,omitempty"`
		Name string `json:"name,omitempty"`
		ID   int64  `json:"id,omitempty"`
	} `json:"contact,omitempty"`
	Watched bool `json:"watched,omitempty"`
}

const contactCollectionType = "application/vnd.ccp.eve.ContactCollection-v1"

type ContactCollectionV1 struct {
	*AuthenticatedClient
	crestPagedFrame
	Items []struct {
		Standing  float64
		Character characterReference
		Contact   struct {
			Href string
			Name string
			ID   int64
		}
		Href        string
		ContactType string
		Watched     bool
		Blocked     bool
	}
}

func (c *AuthenticatedClient) SetContact(id int64, ref string, standing float64) error {
	if err := c.validateClient(); err != nil {
		return err
	}
	contact := ContactCreateV1{Standing: standing}
	contact.Contact.ID = id
	contact.Contact.Href = ref

	url := fmt.Sprintf(c.base.CREST+"characters/%d/contacts/", c.character.CharacterID)
	_, err := c.doJSON("POST", url, contact, nil, contactCreateType)
	if err != nil {
		return err
	}

	return nil
}

func (c *AuthenticatedClient) DeleteContact(id int64, ref string) error {
	if err := c.validateClient(); err != nil {
		return err
	}
	contact := ContactCreateV1{}
	contact.Contact.ID = id
	contact.Contact.Href = ref

	url := fmt.Sprintf(c.base.CREST+"characters/%d/contacts/%d/", c.character.CharacterID, id)
	ret := make(map[string]interface{})
	_, err := c.doJSON("DELETE", url, contact, &ret, contactCreateType)
	if err != nil {
		return err
	}

	return nil
}

func (c *AuthenticatedClient) GetContacts() (*ContactCollectionV1, error) {
	if err := c.validateClient(); err != nil {
		return nil, err
	}
	ret := &ContactCollectionV1{AuthenticatedClient: c}

	url := fmt.Sprintf(c.base.CREST+"characters/%d/contacts/", c.character.CharacterID)

	res, err := c.doJSON("GET", url, nil, ret, contactCollectionType)
	if err != nil {
		return nil, err
	}

	ret.getFrameInfo(url, res)
	return ret, nil
}

func (c *ContactCollectionV1) NextPage() (*ContactCollectionV1, error) {
	w := &ContactCollectionV1{AuthenticatedClient: c.AuthenticatedClient}
	if c.Next.HRef == "" {
		return nil, nil
	}
	res, err := c.doJSON("GET", c.Next.HRef, nil, w, contactCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}
