package main

import (
	"flag"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/dig"
	"github.com/wrfly/web-dns/route"
)

func main() {
	debug := false
	port := *flag.Int("port", 8080, "server port to listen")
	dns := *flag.String("dns", "8.8.8.8:53,8.8.4.4:53", "dns server")
	blacklist := *flag.String("blacklist", "8.8.8.8", "black list")
	flag.BoolVar(&debug, "d", false, "debug switch")
	flag.Parse()

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	digger := dig.New("mem", strings.Split(dns, ","))
	r := route.New(digger, port, strings.Split(blacklist, ","))
	r.Serve()

}
