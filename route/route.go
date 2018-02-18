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
	logrus.Info("create new router")
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
	logrus.Info("start to serve")

	r.e.GET("/:domain", func(c *gin.Context) {
		domain := c.Param("domain")
		logrus.Debugf("got domain: %s", domain)
		if domain == "ping" {
			c.String(200, "pong")
			return
		}
		hosts, err := r.d.Dig(domain, "A")
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
		}
		c.String(http.StatusOK, fmt.Sprintf("%s", hosts))
	})

	r.e.GET("/:domain/:type", func(c *gin.Context) {
		domain := c.Param("domain")
		typ := c.Param("type")
		logrus.Debugf("got domain: %s type: %s", domain, typ)
		if domain == "ping" {
			c.String(200, "pong")
			return
		}

		hosts, err := r.d.Dig(domain, typ)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			c.Abort()
		}
		c.String(http.StatusOK, fmt.Sprintf("%s", hosts))
	})

	// r.e.GET("/json/:domain", func(c *gin.Context) {
	// 	hosts, err := r.d.Dig(c.Param("domain"), "A")
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"err": err.Error(),
	// 		})
	// 		c.Abort()
	// 	}
	// 	c.JSON(http.StatusOK, hosts)
	// })

	// r.e.GET("/json/:domain/:type", func(c *gin.Context) {
	// 	hosts, err := r.d.Dig(c.Param("domain"), c.Param("type"))
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"err": err.Error(),
	// 		})
	// 		c.Abort()
	// 	}
	// 	c.JSON(http.StatusOK, hosts)
	// })

	r.e.Run(fmt.Sprintf(":%d", r.port)) // listen and serve
}
