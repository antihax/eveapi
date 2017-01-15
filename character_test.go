package eveapi

import (
	"net/http"
	"testing"
)

func TestCharacter(t *testing.T) {
	client := &http.Client{}
	r := NewEVEAPIClient(client)
	c, err := r.CharacterV4ByID(1331768660)
	if err != nil {
		t.Errorf("Error getting character %v", err)
	}

	if c.ID != 1331768660 {
		t.Errorf("Character ID does not match the request")
	}
}
