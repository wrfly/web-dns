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
		clientIP := c.GetHeader("X-Forwarded-For")
		if clientIP == "" {
			clientIP = c.ClientIP()
		}
		for _, blk := range blacklist {
			if clientIP == blk {
				c.String(http.StatusForbidden, "403: blacklist")
				return
			}
		}
	}
}

func ratelimitHandler(rate int) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.GetHeader("X-Forwarded-For")
		if clientIP == "" {
			clientIP = c.ClientIP()
		}
		// a client can only query once at one time
		if _, exist := tbs.mC[clientIP]; !exist {
			tbs.mC[clientIP] = &sync.Mutex{}
			// client rate limit
			tbs.tb[clientIP] = ratelimit.New(rate, time.Minute)
			tbs.tbStartAt[clientIP] = time.Now().Second() + 60
		}
		tbs.mC[clientIP].Lock()
		defer tbs.mC[clientIP].Unlock()

		if tbs.tb[clientIP].Limit() {
			tbStartAt := tbs.tbStartAt[clientIP]
			c.String(http.StatusForbidden, "403: ratelimit, try again in %vs\n",
				tbStartAt-time.Now().Second())
			c.Abort()
		}
		return
	}
}

func metricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.GetHeader("X-Forwarded-For")
		if clientIP == "" {
			clientIP = c.ClientIP()
		}
		clientCounter.WithLabelValues(clientIP).Add(1)
	}
}
