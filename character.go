package eveapi

import "fmt"

// CharacterInfo returned data from XML API
type CharacterInfo struct {
	xmlAPIFrame
	CharacterID   int64  `xml:"result>characterID"`
	CharacterName string `xml:"result>characterName"`

	BloodlineID   int64  `xml:"result>bloodlineID"`
	BloodlineName string `xml:"result>bloodline"`

	AncestryID   int64  `xml:"result>ancestryID"`
	AncestryName string `xml:"result>ancestry"`

	CorporationID   int64  `xml:"result>corporationID"`
	CorporationName string `xml:"result>corporation"`

	AllianceID   int64  `xml:"result>allianceID"`
	AllianceName string `xml:"result>alliance"`

	Race string `xml:"result>race"`

	SecurityStatus float64 `xml:"result>securityStatus"`

	EmploymentHistory []struct {
		RecordID        int64      `xml:"recordID,attr"`
		CorporationID   int64      `xml:"corporationID,attr"`
		CorporationName string     `xml:"corporationName,attr"`
		StartDate       EVEXMLTime `xml:"startDate,attr"`
	} `xml:"result>rowset>row"`
}

// GetCharacterInfo queries the XML API for a given characterID.
func (c *AnonymousClient) GetCharacterInfo(characterID int64) (*CharacterInfo, error) {
	w := &CharacterInfo{}

	url := c.base.XML + fmt.Sprintf("eve/CharacterInfo.xml.aspx?characterID=%d", characterID)
	_, err := c.doXML("GET", url, nil, w)
	if err != nil {
		return nil, err
	}
	return w, nil
}

const characterType = "application/vnd.ccp.eve.Character-v4"

type CharacterV4 struct {
	*AnonymousClient
	crestSimpleFrame
	idHref

	Race        idHref
	BloodLine   idHref
	Name        string
	Description string
	Gender      int64
	Corporation entityReference

	Fittings      simpleHref
	Contacts      simpleHref
	Opportunities simpleHref
	Location      simpleHref
	LoyaltyPoints simpleHref

	UI struct {
		SetWaypoints      simpleHref
		ShowContract      simpleHref
		ShowOwnerDetails  simpleHref
		ShowMarketDetails simpleHref
		ShowNewMailWindow simpleHref
	}

	Portrait imageList
}

func (c *AnonymousClient) Character(href string) (*CharacterV4, error) {
	w := &CharacterV4{AnonymousClient: c}
	res, err := c.doJSON("GET", href, nil, w, characterType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(res)
	return w, nil
}

func (c *AnonymousClient) CharacterByID(id int64) (*CharacterV4, error) {
	href := c.base.CREST + fmt.Sprintf("characters/%d/", id)
	return c.Character(href)
}
