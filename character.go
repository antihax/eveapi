package eveapi

import "fmt"

type CharacterInfo struct {
	xmlApiFrame
	CharacterID   int    `xml:"result>characterID"`
	CharacterName string `xml:"result>characterName"`

	EmploymentHistory []struct {
	} `xml:"result>rows"`
}

func (c *AnonymousClient) GetCharacterInfo(characterID int) (*CharacterInfo, error) {
	w := &CharacterInfo{}
	url := c.base.XML + fmt.Sprintf("eve/CharacterInfo.xml.aspx?characterID=%d", characterID)

	res, err := c.httpClient.Get(url)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	err = decodeXML(res, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}
