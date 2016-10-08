package eveapi

import (
	"bytes"
	"crypto/sha1"
	"encoding/xml"
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

type attacker struct {
	AllianceID      int64   `xml:"allianceID,attr"`
	AllianceName    string  `xml:"allianceName,attr"`
	CharacterID     int64   `xml:"characterID,attr"`
	CharacterName   string  `xml:"characterName,attr"`
	CorporationID   int64   `xml:"corporationID,attr"`
	CorporationName string  `xml:"corporationName,attr"`
	DamageDone      int64   `xml:"damageDone,attr"`
	FactionID       int64   `xml:"factionID,attr"`
	FactionName     string  `xml:"factionName,attr"`
	FinalBlow       bool    `xml:"finalBlow,attr"`
	SecurityStatus  float64 `xml:"securityStatus,attr"`
	ShipTypeID      int64   `xml:"shipTypeID,attr"`
	WeaponTypeID    int64   `xml:"weaponTypeID,attr"`
}

type item struct {
	QtyDestroyed int64 `xml:"qtyDestroyed,attr"`
	QtyDropped   int64 `xml:"qtyDropped,attr"`
	TypeID       int64 `xml:"typeID,attr"`
	Flag         int64 `xml:"flag,attr"`
	Singleton    int64 `xml:"singleton,attr"`
}

type KillMailsXML struct {
	xmlAPIFrame

	Kills []struct {
		// Generic kill information
		KillID        int64 `xml:"killID,attr"`
		Hash          string
		SolarSystemID int64      `xml:"solarSystemID,attr"`
		MoonID        int64      `xml:"moonID,attr"`
		KillTime      EVEXMLTime `xml:"killTime,attr"`

		// Victim Information
		Victim struct {
			AllianceID      int64   `xml:"allianceID,attr"`
			AllianceName    string  `xml:"allianceName,attr"`
			CharacterID     int64   `xml:"characterID,attr"`
			CharacterName   string  `xml:"characterName,attr"`
			CorporationID   int64   `xml:"corporationID,attr"`
			CorporationName string  `xml:"corporationName,attr"`
			DamageTaken     int64   `xml:"damageTaken,attr"`
			FactionID       int64   `xml:"factionID,attr"`
			FactionName     string  `xml:"factionName,attr"`
			ShipTypeID      int64   `xml:"shipTypeID,attr"`
			X               float64 `xml:"x,attr"`
			Y               float64 `xml:"y,attr"`
			Z               float64 `xml:"z,attr"`
		} `xml:"victim"`
		Attackers         []attacker `xml:"-"`
		Items             []item     `xml:"-"`
		RawAttackersItems []byte     `xml:",innerxml" json:"-"`
	} `xml:"result>rowset>row"`
}

func (c *AuthenticatedClient) KillMailsXML(characterID int64, fromID int64, rowCount int64) (*KillMailsXML, error) {
	token, err := c.tokenSource.Token()
	if err != nil {
		return nil, err
	}

	v := &KillMailsXML{}

	url := c.base.XML + fmt.Sprintf("char/KillMails.xml.aspx?characterID=%d&accessToken=%s", characterID, token.AccessToken)
	_, err = c.doXML("GET", url, nil, v)
	if err != nil {
		return nil, err
	}

	// Loop through all the kills
	for i, x := range v.Kills {
		// Decode items and attackers.
		v.Kills[i].Attackers, v.Kills[i].Items = decodeAttackerAndItems(x.RawAttackersItems)
		v.Kills[i].RawAttackersItems = []byte{}

		// Find the killing blow
		for _, y := range x.Attackers {
			if y.FinalBlow == true {
				// Compute the kill hash
				v.Kills[i].Hash = GenerateKillMailHash(x.Victim.CharacterID, y.CharacterID, x.Victim.ShipTypeID, x.KillTime.UTC())
				break
			}
		}
	}
	return v, nil
}

// http://stackoverflow.com/questions/39928674/go-unmarshal-nested-arrays-with-identical-name-but-separate-elements
func decodeAttackerAndItems(data []byte) ([]attacker, []item) {
	xmlReader := bytes.NewReader(data)
	decoder := xml.NewDecoder(xmlReader)

	const (
		unknown int = iota
		attackers
		items
	)
	rowset := unknown

	attackerList := []attacker{}
	itemList := []item{}

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "rowset" {
				rowset = unknown
				for _, attr := range se.Attr {
					if attr.Name.Local == "name" {
						if attr.Value == "attackers" {
							rowset = attackers
							break
						} else if attr.Value == "items" {
							rowset = items
							break
						}
					}
				}
			} else if se.Name.Local == "row" {
				switch rowset {
				case attackers:
					a := attacker{}
					if err := decoder.DecodeElement(&a, &se); err == nil {
						attackerList = append(attackerList, a)
					}
				case items:
					it := item{}
					if err := decoder.DecodeElement(&it, &se); err == nil {
						itemList = append(itemList, it)
					}
				}
			}
		}
	}

	return attackerList, itemList
}

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
