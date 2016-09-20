package eveapi

import (
	"fmt"
	"strconv"
	"strings"
)

const marketOrderCollectionSlimType = "application/vnd.ccp.eve.MarketOrderCollectionSlim-v1"

type MarketOrderCollectionSlimV1 struct {
	*AnonymousClient
	crestPagedFrame

	Items []struct {
		Buy           bool
		Issued        EVETime
		Price         float64
		VolumeEntered int64
		MinVolume     int64
		Volume        int64
		Range         string
		Duration      int64
		ID            int64
		Type          int64
		StationID     int64
	}

	RegionID int64 // We wil back fill this for convienence.
}

func (c *AnonymousClient) MarketOrdersSlim(regionID int64, page int) (*MarketOrderCollectionSlimV1, error) {
	w := &MarketOrderCollectionSlimV1{AnonymousClient: c}
	url := c.base.CREST + fmt.Sprintf("market/%d/orders/all/?page=%d", regionID, page)
	res, err := c.doJSON("GET", url, nil, w, marketOrderCollectionSlimType)
	if err != nil {
		return nil, err
	}

	w.RegionID = regionID
	w.getFrameInfo(url, res)
	return w, nil
}

func (c *MarketOrderCollectionSlimV1) NextPage() (*MarketOrderCollectionSlimV1, error) {
	w := &MarketOrderCollectionSlimV1{AnonymousClient: c.AnonymousClient}
	if c.Next.HRef == "" {
		return nil, nil
	}

	res, err := c.doJSON("GET", c.Next.HRef, nil, w, marketOrderCollectionSlimType)
	if err != nil {
		return nil, err
	}

	// [TODO] Improve on this.
	regionID, err := strconv.ParseInt(strings.Split(c.Next.HRef, "/")[4], 10, 64)
	if err != nil {
		return nil, err
	}

	w.RegionID = regionID
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}

const marketTypeHistoryCollectionType = "application/vnd.ccp.eve.MarketTypeHistoryCollection-v1"

type MarketTypeHistoryCollectionV1 struct {
	*AnonymousClient
	crestPagedFrame

	Items []struct {
		OrderCount int64
		LowPrice   float64
		HighPrice  float64
		AvgPrice   float64
		Volume     int64
		Date       string
	}

	RegionID int64 // We wil back fill this for convienence.
	TypeID   int64 // We wil back fill this for convienence.
}

func (c *AnonymousClient) MarketTypeHistoryByID(regionID int64, typeID int64) (*MarketTypeHistoryCollectionV1, error) {
	url := c.base.CREST + fmt.Sprintf("market/%d/history/?type="+c.base.CREST+"inventory/types/%d/", regionID, typeID)
	return c.MarketTypeHistory(url)
}

func (c *AnonymousClient) MarketTypeHistory(url string) (*MarketTypeHistoryCollectionV1, error) {
	w := &MarketTypeHistoryCollectionV1{AnonymousClient: c}

	res, err := c.doJSON("GET", url, nil, w, marketTypeHistoryCollectionType)
	if err != nil {
		return nil, err
	}

	// [TODO] Improve on this.
	regionID, err := strconv.ParseInt(strings.Split(url, "/")[4], 10, 64)
	if err != nil {
		return nil, err
	}
	typeID, err := strconv.ParseInt(strings.Split(url, "/")[11], 10, 64)
	if err != nil {
		return nil, err
	}
	w.RegionID = regionID
	w.TypeID = typeID

	w.getFrameInfo(url, res)
	return w, nil
}

func (c *MarketTypeHistoryCollectionV1) NextPage() (*MarketTypeHistoryCollectionV1, error) {
	if c.Next.HRef == "" {
		return nil, nil
	}

	return c.MarketTypeHistory(c.Next.HRef)
}
