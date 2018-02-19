package main

import (
	"flag"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/config"
	"github.com/wrfly/web-dns/dig"
	"github.com/wrfly/web-dns/route"
)

var (
	debug           bool
	dnsStringList   string
	blackStringList string
	timeOut         int
)

func main() {
	conf := config.Config{}
	flag.IntVar(&conf.Port, "port", 8080, "server port to listen")
	flag.StringVar(&dnsStringList, "dns", "8.8.8.8:53,8.8.4.4:53", "dns server")
	flag.StringVar(&blackStringList, "blacklist", "8.8.8.8", "black list of clients")
	flag.BoolVar(&debug, "d", false, "debug switch")
	flag.IntVar(&timeOut, "t", 100, "dig timeout (millisecond)")
	flag.Parse()

	conf.DNS = strings.Split(dnsStringList, ",")
	conf.BLK = strings.Split(blackStringList, ",")
	conf.Timeout = time.Millisecond * time.Duration(timeOut)

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	digger := dig.New("mem", conf.DNS, conf.Timeout)
	r := route.New(digger, conf.Port, conf.BLK)
	r.Serve()

}
