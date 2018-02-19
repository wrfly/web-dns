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
				c.String(http.StatusForbidden, "403")
				return
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

func (r *Router) stringResp(c *gin.Context, domain, typ string) {
	if domain == "ping" {
		c.String(200, "pong")
		return
	}
	logrus.Debugf("got domain: %s type: %s", domain, typ)
	hosts, err := r.d.Dig(c.Request.Context(), domain, typ)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	if len(hosts) > 0 {
		c.String(http.StatusOK, hosts[0])
	} else {
		c.String(http.StatusOK, "0.0.0.0")
	}
}

func (r *Router) jsonResp(c *gin.Context, domain, typ string) {
	logrus.Debugf("dig domain: %s type: %s", domain, typ)
	answer := r.d.DigJson(c.Request.Context(), domain, typ)
	logrus.Debugf("got answer: %v", answer)
	code := http.StatusOK
	if answer.Err != nil {
		code = http.StatusInternalServerError
	}
	c.JSON(code, answer)
}

func (r *Router) Serve() {
	logrus.Info("start to serve")

	r.e.GET("/:domain", func(c *gin.Context) {
		domain := c.Param("domain")
		r.stringResp(c, domain, "A")
	})

	r.e.GET("/:domain/:type", func(c *gin.Context) {
		domain := c.Param("domain")
		typ := c.Param("type")
		if typ != "json" {
			r.stringResp(c, domain, typ)
		} else { // json response
			r.jsonResp(c, domain, "A")
		}
	})

	r.e.GET("/:domain/:type/json", func(c *gin.Context) {
		r.jsonResp(c, c.Param("domain"), c.Param("type"))
	})

	r.e.Run(fmt.Sprintf(":%d", r.port)) // listen and serve
}
