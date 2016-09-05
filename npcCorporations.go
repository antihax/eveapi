package eveapi

import "fmt"

const npcCorporationsCollectionType = "application/vnd.ccp.eve.NPCCorporationsCollection-v1"

type NPCCorporationsCollectionV1 struct {
	*AnonymousClient
	crestPagedFrame

	Items []struct {
		itemReference
		Description  string
		Headquarters itemReference
		LoyaltyStore struct {
			Href string
		}
		Ticker string
	}
}

func (c *AnonymousClient) NPCCorporations(page int64) (*NPCCorporationsCollectionV1, error) {
	ret := &NPCCorporationsCollectionV1{AnonymousClient: c}
	url := c.base.CREST + fmt.Sprintf("corporations/npccorps/?page=%d", page)

	res, err := c.doJSON("GET", url, nil, ret, npcCorporationsCollectionType)
	if err != nil {
		return nil, err
	}

	ret.getFrameInfo(url, res)
	return ret, nil
}

func (c *NPCCorporationsCollectionV1) NextPage() (*NPCCorporationsCollectionV1, error) {
	w := &NPCCorporationsCollectionV1{AnonymousClient: c.AnonymousClient}
	if c.Next.HRef == "" {
		return nil, nil
	}
	res, err := c.doJSON("GET", c.Next.HRef, nil, w, npcCorporationsCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}

func (c *NPCCorporationsCollectionV1) PreviousPage() (*NPCCorporationsCollectionV1, error) {
	w := &NPCCorporationsCollectionV1{AnonymousClient: c.AnonymousClient}
	if c.Previous.HRef == "" {
		return nil, nil
	}
	res, err := c.doJSON("GET", c.Previous.HRef, nil, w, npcCorporationsCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}
