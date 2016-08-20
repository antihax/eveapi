package eveapi

import "regexp"

// CCP basic XML Frame
type xmlAPIFrame struct {
	Version     int        `xml:"eveapi>version"`
	CurrentTime EVEXMLTime `xml:"currentTime"`
	CachedUntil EVEXMLTime `xml:"cachedUntil"`
}

// IsValidVCode validates a vCode for the XML API meets basic requirements
func IsValidVCode(vc string) bool {
	if m, _ := regexp.MatchString("^[a-zA-Z0-9]+$", vc); !m {
		return false
	}

	if len(vc) != 64 {
		return false
	}

	return true
}
