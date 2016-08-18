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

type Contacts struct {
	*AuthenticatedClient
	crestPagedFrame
	Items []struct {
		Standing  float64
		Character struct {
			Name        string
			Corporation struct {
				Name  string
				IsNPC bool
				Href  string

				Logo struct {
					I32x32 struct {
						Href string
					}
					I64x64 struct {
						Href string
					}
					I128x128 struct {
						Href string
					}
					I256x256 struct {
						Href string
					}
				}
				ID int64
			}
			Alliance struct {
				Name  string
				IsNPC bool
				Href  string

				Logo struct {
					I32x32 struct {
						Href string
					}
					I64x64 struct {
						Href string
					}
					I128x128 struct {
						Href string
					}
					I256x256 struct {
						Href string
					}
				}
				ID int64
			}
			IsNPC     bool
			Href      string
			Capsuleer struct {
				Href string
			}
			Portrait struct {
				I32x32 struct {
					Href string
				}
				I64x64 struct {
					Href string
				}
				I128x128 struct {
					Href string
				}
				I256x256 struct {
					Href string
				}
			}
			ID int64
		}
		Contact struct {
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
	contact := Contact{Standing: standing}
	contact.Contact.ID = id
	contact.Contact.Href = ref

	url := fmt.Sprintf(c.base.CREST+"characters/%d/contacts/", c.character.CharacterID)
	_, err := c.doJSON("POST", url, contact, nil)
	if err != nil {
		return err
	}

	return nil
}

func (c *AuthenticatedClient) DeleteContact(id int64, ref string) error {
	if err := c.validateClient(); err != nil {
		return err
	}
	contact := Contact{}
	contact.Contact.ID = id
	contact.Contact.Href = ref

	url := fmt.Sprintf(c.base.CREST+"characters/%d/contacts/%d/", c.character.CharacterID, id)
	ret := make(map[string]interface{})
	_, err := c.doJSON("DELETE", url, contact, &ret)
	if err != nil {
		return err
	}

	return nil
}

func (c *AuthenticatedClient) GetContacts() (*Contacts, error) {
	if err := c.validateClient(); err != nil {
		return nil, err
	}
	ret := &Contacts{AuthenticatedClient: c}

	url := fmt.Sprintf(c.base.CREST+"characters/%d/contacts/", c.character.CharacterID)

	res, err := c.doJSON("GET", url, nil, ret)
	if err != nil {
		return nil, err
	}

	ret.getFrameInfo(url, res)
	return ret, nil
}

func (c *Contacts) NextPage() (*Contacts, error) {
	w := &Contacts{AuthenticatedClient: c.AuthenticatedClient}
	if c.Next.HRef == "" {
		return nil, nil
	}
	res, err := c.doJSON("GET", c.Next.HRef, nil, w)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}
