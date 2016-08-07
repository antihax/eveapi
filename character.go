package eveapi

import "fmt"

type CharacterInfo struct {
	xmlApiFrame
	CharacterID   int    `xml:"result>characterID"`
	CharacterName string `xml:"result>characterName"`

	BloodlineID   int    `xml:"result>bloodlineID"`
	BloodlineName string `xml:"result>bloodline"`

	AncestryID   int    `xml:"result>ancestryID"`
	AncestryName string `xml:"result>ancestry"`

	CorporationID   int    `xml:"result>corporationID"`
	CorporationName string `xml:"result>corporation"`

	AllianceID   int    `xml:"result>allianceID"`
	AllianceName string `xml:"result>alliance"`

	Race string `xml:"result>race"`

	SecurityStatus float64 `xml:"result>securityStatus"`

	EmploymentHistory []struct {
		RecordID        int        `xml:"recordID,attr"`
		CorporationID   int        `xml:"corporationID,attr"`
		CorporationName string     `xml:"corporationName,attr"`
		StartDate       EVEXMLTime `xml:"startDate,attr"`
	} `xml:"result>rowset>row"`
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
