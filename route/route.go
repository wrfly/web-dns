package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"

	"github.com/wrfly/web-dns/config"
	"github.com/wrfly/web-dns/dig"
)

type Router struct {
	port      int
	debugPort int
	debug     bool
	blacklist []string
	d         dig.Digger
	e         *gin.Engine
}

func New(digger dig.Digger, conf config.ServerConfig) *Router {
	logrus.Info("create new router")
	engine := gin.New()
	engine.Use(blacklistHandler(conf.BLK))
	engine.Use(ratelimitHandler(conf.Rate))

	debug := false
	if conf.DebugPort > 0 {
		debug = true
		registerPrometheus()
		engine.Use(metricsHandler())
		go serveMetricsAndDebug(conf.DebugPort)
	}

	return &Router{
		port:      conf.Port,
		debugPort: conf.DebugPort,
		blacklist: conf.BLK,
		debug:     debug,
		d:         digger,
		e:         engine,
	}
}

func (r *Router) stringResp(c *gin.Context, domain, typ string) {
	if domain == "ping" {
		c.String(200, "pong")
		return
	}
	if r.debug {
		domainCounter.WithLabelValues(domain, typ).Add(1)
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

func serveMetricsAndDebug(debugPort int) {
	logrus.Info("start to serve metrics and debug")
	s := http.NewServeMux()
	s.Handle("/metrics", prometheus.Handler())
	s.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pprof.Index(w, r)
	})

	addr := fmt.Sprintf("127.0.0.1:%d", debugPort)
	http.ListenAndServe(addr, s)

}
