package eveapi

import "fmt"

const killmailType = "application/vnd.ccp.eve.Killmail-v1+json"

type KillmailV1 struct {
	*AnonymousClient
	crestSimpleFrame
	Hash string

	SolarSystem itemReference
	KillID      int64
	KillTime    EVEKillmailTime
	Attackers   []struct {
		ShipType       itemReference
		Corporation    itemReference
		Character      itemReference
		Alliance       itemReference
		WeaponType     itemReference
		FinalBlow      bool
		SecurityStatus float64
		DamageDone     int64
	}
	AttackerCount int64
	Victim        struct {
		DamageTaken int64
		Items       []struct {
			Singleton         int64
			ItemType          itemReference
			Flag              int64
			QuantityDestroyed int64
			QuantityDropped   int64
		}

		Character   itemReference
		ShipType    itemReference
		Corporation itemReference
		Alliance    itemReference
		Position    struct {
			X float64
			Y float64
			Z float64
		}
	}
	War struct {
		Href string
		ID   int64
	}
}

// Killmail pulls killmail information
func (c *AnonymousClient) Killmail(href string) (*KillmailV1, error) {
	w := &KillmailV1{AnonymousClient: c}
	res, err := c.doJSON("GET", href, nil, w, warType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(res)
	return w, nil
}

// KillmailByID requires the killmail ID and hash string for validation and collects the killmail.
func (c *AnonymousClient) KillmailByID(id int, hash string) (*KillmailV1, error) {
	w := &KillmailV1{AnonymousClient: c}
	url := c.base.CREST + fmt.Sprintf("killmails/%d/%s/", id, hash)
	res, err := c.doJSON("GET", url, nil, w, warType)
	if err != nil {
		return nil, err
	}
	w.getFrameInfo(res)
	return w, nil
}
