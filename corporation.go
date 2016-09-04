package eveapi

import "fmt"

// CharacterInfo returned data from XML API
type CorporationSheet struct {
	xmlAPIFrame
	CorporationID   int64  `xml:"result>corporationID"`
	CorporationName string `xml:"result>corporationName"`
	Ticker          string `xml:"result>ricker"`
	CEOID           int64  `xml:"result>ceoID"`
	CEOName         string `xml:"result>ceoName"`
	StationID       int64  `xml:"result>stationID"`
	StationName     string `xml:"result>stationName"`
	Description     string `xml:"result>description"`
	AllianceID      int64  `xml:"result>allianceID"`
	AllianceName    string `xml:"result>allianceName"`
	FactionID       int64  `xml:"result>factionID"`
	URL             string `xml:"result>url"`
	MemberCount     int64  `xml:"result>memberCount"`
	Shares          int64  `xml:"result>shares"`
	Logo            struct {
		GraphicID int64 `xml:"grapicID,attr"`
		Shape1    int64 `xml:"shape1,attr"`
		Shape2    int64 `xml:"shape2,attr"`
		Shape3    int64 `xml:"shape3,attr"`
		Color1    int64 `xml:"color1,attr"`
		Color2    int64 `xml:"color2,attr"`
		Color3    int64 `xml:"color3,attr"`
	} `xml:"result>logo"`
}

const loyaltyStoreOffersCollectionType = "application/vnd.ccp.eve.LoyaltyStoreOffersCollection-v1"

type LoyaltyStoreOffersCollectionV1 struct {
	*AnonymousClient
	crestPagedFrame

	Items []struct {
		IskCost int64
		LpCost  int64
		Item    struct {
			itemReference
			RequiredItems []struct {
				Item     itemReference
				Quantity int64
			}
		}
	}
}

func (c *AnonymousClient) LoyaltyPointStore(corporationID int64) (*LoyaltyStoreOffersCollectionV1, error) {
	ret := &LoyaltyStoreOffersCollectionV1{AnonymousClient: c}
	url := fmt.Sprintf(c.base.CREST+"corporations/%d/loyaltystore/", corporationID)

	res, err := c.doJSON("GET", url, nil, ret, loyaltyStoreOffersCollectionType)
	if err != nil {
		return nil, err
	}

	ret.getFrameInfo(url, res)
	return ret, nil
}

func (c *LoyaltyStoreOffersCollectionV1) NextPage() (*LoyaltyStoreOffersCollectionV1, error) {
	w := &LoyaltyStoreOffersCollectionV1{AnonymousClient: c.AnonymousClient}
	if c.Next.HRef == "" {
		return nil, nil
	}
	res, err := c.doJSON("GET", c.Next.HRef, nil, w, loyaltyStoreOffersCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}

func (c *LoyaltyStoreOffersCollectionV1) PreviousPage() (*LoyaltyStoreOffersCollectionV1, error) {
	w := &LoyaltyStoreOffersCollectionV1{AnonymousClient: c.AnonymousClient}
	if c.Previous.HRef == "" {
		return nil, nil
	}
	res, err := c.doJSON("GET", c.Previous.HRef, nil, w, loyaltyStoreOffersCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}

// GetCharacterInfo queries the XML API for a given characterID.
func (c *AnonymousClient) GetCorporationPublicSheet(corporationID int64) (*CorporationSheet, error) {
	w := &CorporationSheet{}

	url := c.base.XML + fmt.Sprintf("corp/CorporationSheet.xml.aspx?corporationID=%d", corporationID)
	_, err := c.doXML("GET", url, nil, w)
	if err != nil {
		return nil, err
	}
	return w, nil
}
