package route

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/wrfly/web-dns/config"
	"github.com/wrfly/web-dns/dig"
)

type Router struct {
	port      int
	debugPort int
	blacklist []string
	d         dig.Digger
	e         *gin.Engine
}

func New(digger dig.Digger, conf config.ServerConfig) *Router {
	logrus.Info("create new router")
	engine := gin.New()
	engine.Use(blacklistHandler(conf.BLK))
	engine.Use(ratelimitHandler(conf.Rate))

	return &Router{
		port:      conf.Port,
		debugPort: conf.DebugPort,
		blacklist: conf.BLK,
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
	c.JSON(code, answer.Result)
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

	rs := r.e.Routes()
	paths := []string{}
	for _, ri := range rs {
		paths = append(paths, ri.Path)
	}
	bs, _ := json.MarshalIndent(gin.H{
		"Paths":  paths,
		"Readme": "github.com/wrfly/web-dns",
	}, "", "    ")
	usage := fmt.Sprintf("%s", bs)
	r.e.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, usage)
	})

	r.e.Run(fmt.Sprintf(":%d", r.port)) // listen and serve
}
