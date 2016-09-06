package eveapi

import "regexp"

// CCP basic XML Frame
type xmlAPIFrame struct {
	Version     int        `xml:"eveapi>version"`
	CurrentTime EVEXMLTime `xml:"currentTime"`
	CachedUntil EVEXMLTime `xml:"cachedUntil"`
}

// XMLAPIKey holds an API key for the XML API.
type XMLAPIKey struct {
	VCode string
	KeyID int64
}

// IsValidVCode validates a vCode for the XML API meets basic requirements
func (c *XMLAPIKey) IsValidVCode() bool {
	if m, _ := regexp.MatchString("^[a-zA-Z0-9]+$", c.VCode); !m {
		return false
	}

	if len(c.VCode) != 64 {
		return false
	}

	return true
}
