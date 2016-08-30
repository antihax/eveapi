package eveapi

import "fmt"

// CharacterInfo returned data from XML API
type characterInfo struct {
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
func (c *AnonymousClient) GetCharacterInfo(characterID int64) (*characterInfo, error) {
	w := &characterInfo{}

	url := c.base.XML + fmt.Sprintf("eve/CharacterInfo.xml.aspx?characterID=%d", characterID)
	_, err := c.doXML("GET", url, nil, w)
	if err != nil {
		return nil, err
	}
	return w, nil
}
