package eveapi

import "time"

// CCP's documentation states rate limits are tracked by IP address.
// The throttles provide a burstable rate limit to each component of the API.

// Authenticated SSO Throttle
var authedThrottle = make(chan time.Time, 10) // Burst 20

// Anonymous CREST Throttle
var anonThrottle = make(chan time.Time, 100) // Burst 100
var anonConnectionLimit = make(chan bool, 20)

// XML client Throttle
var xmlThrottle = make(chan time.Time, 20) // Burst 30

func init() {
	// Authenticated SSO client rate limit
	var authedRate = time.Second / 10
	var authedTick = time.NewTicker(authedRate)

	go func() {
		for t := range authedTick.C {
			select {
			case authedThrottle <- t:
			default:
			}
		}
	}()

	// Anonymous CREST client rate limit
	var anonRate = time.Second / 100
	var anonTick = time.NewTicker(anonRate)
	go func() {
		for t := range anonTick.C {
			select {
			case anonThrottle <- t:
			default:
			}
		}
	}()

	// XML client rate limit
	var xmlRate = time.Second / 20
	var xmlTick = time.NewTicker(xmlRate)
	go func() {
		for t := range xmlTick.C {
			select {
			case xmlThrottle <- t:
			default:
			}
		}
	}()
}
