package eveapi

import "fmt"

const loyaltyStoreOffersCollectionV1Type = "application/vnd.ccp.eve.LoyaltyStoreOffersCollection-v1"

type LoyaltyStoreOffersCollectionV1 struct {
	*AnonymousClient
	crestPagedFrame

	Items []struct {
		ID            int64
		AkCost        int64
		IskCost       int64
		LpCost        int64
		Quantity      int64
		Item          itemReference
		RequiredItems []struct {
			Item     itemReference
			Quantity int64
		}
	}
}

func (c *AnonymousClient) LoyaltyPointStoreV1(url string) (*LoyaltyStoreOffersCollectionV1, error) {
	ret := &LoyaltyStoreOffersCollectionV1{AnonymousClient: c}

	res, err := c.doJSON("GET", url, nil, ret, loyaltyStoreOffersCollectionV1Type)
	if err != nil {
		return nil, err
	}

	ret.getFrameInfo(url, res)
	return ret, nil
}

func (c *AnonymousClient) LoyaltyPointStoreV1ByID(corporationID int64) (*LoyaltyStoreOffersCollectionV1, error) {
	url := fmt.Sprintf(c.base.CREST+"corporations/%d/loyaltystore/", corporationID)
	return c.LoyaltyPointStoreV1(url)
}

func (c *LoyaltyStoreOffersCollectionV1) NextPage() (*LoyaltyStoreOffersCollectionV1, error) {
	w := &LoyaltyStoreOffersCollectionV1{AnonymousClient: c.AnonymousClient}
	if c.Next.HRef == "" {
		return nil, nil
	}

	res, err := c.doJSON("GET", c.Next.HRef, nil, w, loyaltyStoreOffersCollectionV1Type)
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
	res, err := c.doJSON("GET", c.Previous.HRef, nil, w, loyaltyStoreOffersCollectionV1Type)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}
