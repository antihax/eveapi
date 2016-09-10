package eveapi

import "time"

type rateLimiter struct {
	throttle  chan time.Time
	rate      time.Duration
	ticker    *time.Ticker
	burstRate int
}

func newRateLimiter(requestsPerSecond int, burstRate int) *rateLimiter {
	c := &rateLimiter{
		throttle:  make(chan time.Time, burstRate),
		rate:      time.Second / (time.Duration)(requestsPerSecond),
		burstRate: burstRate,
	}
	c.startRatelimiter()
	return c
}

func (c *rateLimiter) startRatelimiter() {
	// Create the timed limiter
	c.ticker = time.NewTicker(c.rate)

	// Fill the buffer with the burst tokens
	for i := 0; i < c.burstRate; i++ {
		c.throttle <- time.Now()
	}

	// Start the rate limiter
	go c.tick()
}

func (c *rateLimiter) tick() {
	for t := range c.ticker.C {
		select {
		case c.throttle <- t:
		default:
		}
	}
}

func (c *rateLimiter) stop() {
	c.ticker.Stop()
}

func (c *rateLimiter) throttleRequest() {
	<-c.throttle
}

// Prevent going over 20 connections on anonymous clients
var anonConnectionLimit = make(chan bool, 20)

// CCP's documentation states rate limits are tracked by IP address.
// The throttles provide a burstable rate limit to each component of the API.
var authedThrottle *rateLimiter
var anonThrottle *rateLimiter
var xmlThrottle *rateLimiter

func init() {
	authedThrottle = newRateLimiter(20, 20)
	anonThrottle = newRateLimiter(150, 300)
	xmlThrottle = newRateLimiter(30, 30)
}
