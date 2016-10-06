package eveapi

import "fmt"

const alliancesCollectionType = "application/vnd.ccp.eve.AlliancesCollection-v2"

type AlliancesCollectionV2 struct {
	*AnonymousClient
	crestPagedFrame
	Items []struct {
		ShortName string
		HRef      string
		ID        int64
		Name      string
	}
}

func (c *AnonymousClient) Alliances(page int) (*AlliancesCollectionV2, error) {
	w := &AlliancesCollectionV2{AnonymousClient: c}
	url := c.base.CREST + fmt.Sprintf("alliances/?page=%d", page)
	res, err := c.doJSON("GET", url, nil, w, alliancesCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(url, res)
	return w, nil
}

func (c *AlliancesCollectionV2) NextPage() (*AlliancesCollectionV2, error) {
	w := &AlliancesCollectionV2{AnonymousClient: c.AnonymousClient}
	if c.Next.HRef == "" {
		return nil, nil
	}
	res, err := c.doJSON("GET", c.Next.HRef, nil, w, alliancesCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}

const allianceType = "application/vnd.ccp.eve.Alliance-v1"

type AllianceV1 struct {
	*AnonymousClient
	crestSimpleFrame
	StartDate           EVETime
	CorporationsCount   int64
	Description         string
	ExecutorCorporation entityReference
	CreatorCorporation  entityReference
	URL                 string
	ID                  int64
	Name                string
	ShortName           string
	Deleted             bool
	CreatorCharacter    characterReference
	Corporations        []entityReference
}

func (c *AnonymousClient) Alliance(href string) (*AllianceV1, error) {
	w := &AllianceV1{AnonymousClient: c}
	res, err := c.doJSON("GET", href, nil, w, allianceType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(res)
	return w, nil
}

func (c *AnonymousClient) AllianceByID(id int64) (*AllianceV1, error) {
	href := c.base.CREST + fmt.Sprintf("alliances/%d/", id)
	return c.Alliance(href)
}
