package eveapi

import "fmt"

const allianceType = "application/vnd.ccp.eve.Alliance-v1+json"

type allianceV1 struct {
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

const alliancesCollectionType = "application/vnd.ccp.eve.AlliancesCollection-v2+json"

type alliancesCollectionV2 struct {
	*AnonymousClient
	crestPagedFrame
	Items []struct {
		ShortName string
		HRef      string
		ID        int64
		Name      string
	}
}

func (c *AnonymousClient) Alliances(page int) (*alliancesCollectionV2, error) {
	w := &alliancesCollectionV2{AnonymousClient: c}
	url := c.base.CREST + fmt.Sprintf("alliances/?page=%d", page)
	res, err := c.doJSON("GET", url, nil, w, alliancesCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(url, res)
	return w, nil
}

func (c *alliancesCollectionV2) NextPage() (*alliancesCollectionV2, error) {
	w := &alliancesCollectionV2{AnonymousClient: c.AnonymousClient}
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

func (c *AnonymousClient) Alliance(href string) (*allianceV1, error) {
	w := &allianceV1{AnonymousClient: c}
	res, err := c.doJSON("GET", href, nil, w, allianceType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(res)
	return w, nil
}

func (c *AnonymousClient) AllianceByID(id int64) (*allianceV1, error) {
	w := &allianceV1{AnonymousClient: c}
	href := c.base.CREST + fmt.Sprintf("alliances/%d/", id)
	res, err := c.doJSON("GET", href, nil, w, allianceType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(res)
	return w, nil
}
