package eveapi

import "fmt"

const warsCollectionType = "application/vnd.ccp.eve.WarsCollection-v1+json"

type warsCollectionV1 struct {
	*AnonymousClient
	crestPagedFrame

	Items []struct {
		HRef string
		ID   int
	}
}

const warKillmailsType = "application/vnd.ccp.eve.WarKillmails-v1"

type warKillmailsV1 struct {
	*AnonymousClient
	crestPagedFrame

	Items []struct {
		HRef string
		ID   int
	}
}

const warType = "application/vnd.ccp.eve.War-v1+json"

type warV1 struct {
	*AnonymousClient
	crestSimpleFrame

	TimeFinished  EVETime
	OpenForAllies bool
	TimeStarted   EVETime
	AllyCount     int
	TimeDeclared  EVETime

	Allies []struct {
		HRef string
		ID   int64
		Icon struct {
			HRef string
		}
		Name string
	}
	Aggressor struct {
		ShipsKilled int

		Name string
		HRef string

		Icon struct {
			HRef string
		}
		ID        int64
		IskKilled float64
	}
	Mutual bool

	Killmails string

	Defender struct {
		ShipsKilled int

		Name string
		HRef string

		Icon struct {
			HRef string
		}
		ID        int64
		IskKilled float64
	}
	ID int64
}

func (c *AnonymousClient) Wars(page int) (*warsCollectionV1, error) {
	w := &warsCollectionV1{AnonymousClient: c}
	url := c.base.CREST + fmt.Sprintf("wars/?page=%d", page)
	res, err := c.doJSON("GET", url, nil, w, warsCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(url, res)
	return w, nil
}

func (c *warsCollectionV1) NextPage() (*warsCollectionV1, error) {
	w := &warsCollectionV1{AnonymousClient: c.AnonymousClient}
	if c.Next.HRef == "" {
		return nil, nil
	}
	res, err := c.doJSON("GET", c.Next.HRef, nil, w, warsCollectionType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}

func (c *AnonymousClient) War(href string) (*warV1, error) {
	w := &warV1{AnonymousClient: c}
	res, err := c.doJSON("GET", href, nil, w, warType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(res)
	return w, nil
}

func (c *AnonymousClient) WarByID(id int) (*warV1, error) {
	w := &warV1{AnonymousClient: c}
	url := c.base.CREST + fmt.Sprintf("wars/%d/", id)
	res, err := c.doJSON("GET", url, nil, w, warType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(res)
	return w, nil
}

// GetKillmails provides a list of killmails associated to this war.
func (c *warV1) GetKillmails() (*warKillmailsV1, error) {
	w := &warKillmailsV1{AnonymousClient: c.AnonymousClient}
	res, err := c.doJSON("GET", c.Killmails, nil, w, warKillmailsType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(c.Killmails, res)
	return w, nil
}
