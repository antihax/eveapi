package eveapi

import (
	"fmt"
	"regexp"
)

// CharacterInfo returned data from XML API
type CharacterInfoXML struct {
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
func (c *EVEAPIClient) CharacterInfoXML(characterID int64) (*CharacterInfoXML, error) {
	w := &CharacterInfoXML{}

	url := c.base.XML + fmt.Sprintf("eve/CharacterInfo.xml.aspx?characterID=%d", characterID)
	_, err := c.doXML("GET", url, nil, w, nil)
	if err != nil {
		return nil, err
	}
	return w, nil
}

const characterV4Type = "application/vnd.ccp.eve.Character-v4"

type CharacterV4 struct {
	*EVEAPIClient
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

func (c *EVEAPIClient) CharacterV4(href string) (*CharacterV4, error) {
	w := &CharacterV4{EVEAPIClient: c}
	res, err := c.doJSON("GET", href, nil, w, characterV4Type)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(res)
	return w, nil
}

func (c *EVEAPIClient) CharacterV4ByID(id int64) (*CharacterV4, error) {
	href := c.base.CREST + fmt.Sprintf("characters/%d/", id)
	return c.CharacterV4(href)
}

// https://community.eveonline.com/support/policies/naming-policy-en/
func ValidCharacterName(name string) bool {
	if len(name) > 37 {
		return false
	}
	if m, _ := regexp.MatchString("^[a-zA-Z0-9' -]+$", name); !m {
		return false
	}
	return true
}
