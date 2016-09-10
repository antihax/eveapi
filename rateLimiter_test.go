package eveapi

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	r := newRateLimiter(1, 20)
	count := 0
	start := time.Now()
	for {
		count++
		r.throttleRequest()
		if count == 20 {
			if (int)(time.Now().Sub(start).Seconds()) != 0 {
				t.Errorf("Burst failed") // This should be immediate
			}
		}
		if count == 22 {
			if (int)(time.Now().Sub(start).Seconds()) != 2 { // after two seconds
				t.Errorf("Rate limiter failed to properly limit %d", (int)(time.Now().Sub(start)))
			}
			time.Sleep(time.Second * 2)
		}

		if count == 25 {
			if (int)(time.Now().Sub(start).Seconds()) != 5 { // after five seconds
				t.Errorf("Failed to recover burst tokens %d", (int)(time.Now().Sub(start)))
			}
			break
		}
	}
}
