package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/wrfly/web-dns/dig"
)

type Router struct {
	port      int
	blacklist []string
	d         dig.Digger
	e         *gin.Engine
}

func blacklistHandler(blacklist []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cip := c.GetHeader("X-Forwarded-For")
		if cip == "" {
			cip = c.ClientIP()
		}
		for _, blk := range blacklist {
			if cip == blk {
				c.String(http.StatusForbidden, "emm")
				c.Abort()
			}
		}
	}
}

func New(digger dig.Digger, port int, blacklist []string) *Router {
	engine := gin.New()
	engine.Use(blacklistHandler(blacklist))

	return &Router{
		port:      port,
		blacklist: blacklist,
		d:         digger,
		e:         engine,
	}
}

func (r *Router) Serve() {
	logrus.Debug("server starting")
	r.e.Any("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.e.GET("/:host")

	r.e.Run(fmt.Sprintf(":%d", r.port)) // listen and serve
}
