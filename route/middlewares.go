package route

import (
	"net/http"
	"sync"
	"time"

	"github.com/bsm/ratelimit"
	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	tb        map[string]*ratelimit.RateLimiter
	tbStartAt map[string]int
	mC        map[string]*sync.Mutex
}

var tbs = rateLimiter{
	tb:        make(map[string]*ratelimit.RateLimiter),
	tbStartAt: make(map[string]int),
	mC:        make(map[string]*sync.Mutex),
}

func blacklistHandler(blacklist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cip := c.GetHeader("X-Forwarded-For")
		if cip == "" {
			cip = c.ClientIP()
		}
		for _, blk := range blacklist {
			if cip == blk {
				c.String(http.StatusForbidden, "403: blacklist")
				return
			}
		}
	}
}

func ratelimitHandler(rate int) gin.HandlerFunc {
	return func(c *gin.Context) {
		cip := c.GetHeader("X-Forwarded-For")
		if cip == "" {
			cip = c.ClientIP()
		}
		// a client can only query once at one time
		if _, exist := tbs.mC[cip]; !exist {
			tbs.mC[cip] = &sync.Mutex{}
			// client rate limit
			tbs.tb[cip] = ratelimit.New(rate, time.Minute)
			tbs.tbStartAt[cip] = time.Now().Second() + 60
		}
		tbs.mC[cip].Lock()
		defer tbs.mC[cip].Unlock()

		if tbs.tb[cip].Limit() {
			tbStartAt := tbs.tbStartAt[cip]
			c.String(http.StatusForbidden, "403: ratelimit, try again in %vs\n",
				tbStartAt-time.Now().Second())
			c.Abort()
		}
		return
	}
}
