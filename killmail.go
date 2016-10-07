package eveapi

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const killmailV1Type = "application/vnd.ccp.eve.Killmail-v1"

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
func (c *AnonymousClient) KillmailV1(href string) (*KillmailV1, error) {
	w := &KillmailV1{AnonymousClient: c}
	res, err := c.doJSON("GET", href, nil, w, killmailV1Type)
	if err != nil {
		return nil, err
	}

	// Generate extra frame information
	w.getFrameInfo(res)

	// Add the hash to the structure, pull it from the URL
	w.Hash = strings.Split(w.PageURL, "/")[5]
	return w, nil
}

// KillmailByID requires the killmail ID and hash string for validation and collects the killmail.
func (c *AnonymousClient) KillmailV1ByID(id int, hash string) (*KillmailV1, error) {
	w := &KillmailV1{AnonymousClient: c}
	url := c.base.CREST + fmt.Sprintf("killmails/%d/%s/", id, hash)
	res, err := c.doJSON("GET", url, nil, w, killmailV1Type)
	if err != nil {
		return nil, err
	}

	w.Hash = hash
	w.getFrameInfo(res)
	return w, nil
}

/*
// [TODO]
type KillMailsXML struct {
	xmlAPIFrame
}

func (c *AuthenticatedClient) KillMailsXML(characterID int64) (*KillMailsXML, error) {
	w := &KillMailsXML{}

	url := c.base.XML + fmt.Sprintf("char/KillMails.xml.aspx?characterID=%d", characterID)
	_, err := c.doXML("GET", url, nil, w)
	if err != nil {
		return nil, err
	}
	return w, nil
}
*/

// Generate the killmail hash using source information.
func GenerateKillMailHash(victimID int64, attackerID int64, shipTypeID int64, killTime time.Time) string {
	v := strconv.FormatInt(victimID, 10)
	if victimID == 0 {
		v = "None"
	}
	a := strconv.FormatInt(attackerID, 10)
	if attackerID == 0 {
		a = "None"
	}
	s := strconv.FormatInt(shipTypeID, 10)

	t := strconv.FormatInt(convertTimeToWindow64(killTime), 10)

	h := sha1.New()
	io.WriteString(h, v)
	io.WriteString(h, a)
	io.WriteString(h, s)
	io.WriteString(h, t)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Convert to killmail time for hashing.
func convertTimeToWindow64(t time.Time) int64 {
	return t.Unix()*10000000 + 116444736000000000
}
