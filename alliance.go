package eveapi

import (
	"fmt"
	"time"
)

type Corporation struct {
	Name  string
	isNPC bool
	Href  string
	Logo  struct {
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

const allianceType = "application/vnd.ccp.eve.Alliance-v1+json"

type AllianceV1 struct {
	*AnonymousClient
	crestSimpleFrame
	StartDate           time.Time
	CorporationsCount   int64
	Description         string
	ExecutorCorporation Corporation
	URL                 string
	ID                  int64
	Name                string
	ShortName           string
}

const alliancesCollectionType = "application/vnd.ccp.eve.AlliancesCollection-v2+json"

type AlliancesCollectionV2 struct {
	*AnonymousClient
	crestPagedFrame
	Items []struct {
		ShortName string
		Href      string
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
