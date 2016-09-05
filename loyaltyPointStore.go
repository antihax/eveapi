package eveapi

import "fmt"

const loyaltyStoreOffersCollectionType = "application/vnd.ccp.eve.LoyaltyStoreOffersCollection-v1"

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

func (c *AnonymousClient) LoyaltyPointStoreByID(corporationID int64) (*LoyaltyStoreOffersCollectionV1, error) {
	url := fmt.Sprintf(c.base.CREST+"corporations/%d/loyaltystore/", corporationID)
	return c.LoyaltyPointStore(url)
}

func (c *AnonymousClient) LoyaltyPointStore(url string) (*LoyaltyStoreOffersCollectionV1, error) {
	ret := &LoyaltyStoreOffersCollectionV1{AnonymousClient: c}

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
