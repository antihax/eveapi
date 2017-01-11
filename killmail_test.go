package eveapi

import (
	"net/http"
	"strconv"
	"testing"
	"time"
)

func TestKillmail(t *testing.T) {
	client := &http.Client{}
	r := NewEVEAPIClient(client)
	c, err := r.KillmailV1ByID(40583728, "efd4bf9c4f2aee704d3f9a7f8ae0176a15eba19d")
	if err != nil {
		t.Errorf("Error getting killmail %v", err)
	}

	if c.Hash != "efd4bf9c4f2aee704d3f9a7f8ae0176a15eba19d" {
		t.Errorf("Hash does not match the request")
	}

	c, err = r.KillmailV1("https://crest-tq.eveonline.com/killmails/40583728/efd4bf9c4f2aee704d3f9a7f8ae0176a15eba19d/")
	if err != nil {
		t.Errorf("Error getting killmail %v", err)
	}

	if c.Hash != "efd4bf9c4f2aee704d3f9a7f8ae0176a15eba19d" {
		t.Errorf("Hash does not match the request")
	}
}

func TestKillmailHash(t *testing.T) {
	test := "2014.08.10 01:58:00"
	ti, err := time.Parse(eveKillmailTimeLayout, test)
	if err != nil {
		t.Errorf("Error parsing time %v", err)
	}
	hash := GenerateKillMailHash(0, 93808108, 24646, ti)

	if hash != "efd4bf9c4f2aee704d3f9a7f8ae0176a15eba19d" {
		t.Errorf("Hash does not match")
	}
}

func TestKillmailTime(t *testing.T) {
	test := "2014.08.10 01:58:00"
	ti, err := time.Parse(eveKillmailTimeLayout, test)
	if err != nil {
		t.Errorf("Error parsing time %v", err)
	}
	d := convertTimeToWindow64(ti)
	if d != 130521094800000000 {
		i, err := strconv.ParseInt("1305210948", 10, 64)
		if err != nil {
			panic(err)
		}
		tm := time.Unix(i, 0)
		t.Errorf("Time does not match %d != %d %s %s", 130521094800000000, d, ti.String(), tm.String())
	}
}
