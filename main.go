package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/config"
	"github.com/wrfly/web-dns/dig"
	"github.com/wrfly/web-dns/dig/cache"
	"github.com/wrfly/web-dns/route"
)

var (
	debug           bool
	dnsStringList   string
	blackStringList string
	timeOut         int
)

func main() {
	conf := config.New()
	flag.IntVar(&conf.Server.Port, "port", 8080, "server port to listen")
	flag.StringVar(&dnsStringList, "dns", "8.8.8.8:53,8.8.4.4:53", "dns server")
	flag.StringVar(&blackStringList, "blacklist", "8.8.8.8", "black list of clients")
	flag.BoolVar(&conf.Debug, "d", false, "debug switch")
	flag.IntVar(&timeOut, "t", 100, "dig timeout (millisecond)")
	flag.IntVar(&conf.Server.Rate, "r", 1000, "rate of requests per minute per IP")
	flag.StringVar(&conf.Cacher.CacheType, "cache", "mem", "cache type: mem|redis|bolt")
	flag.StringVar(&conf.Cacher.RedisAddr, "redis", "localhost:6379", "this flag is used for redis cacher")
	flag.Parse()

	conf.Server.BLK = strings.Split(blackStringList, ",")
	conf.Digger.DNS = strings.Split(dnsStringList, ",")
	conf.Digger.Timeout = time.Millisecond * time.Duration(timeOut)

	if conf.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode")
	} else {
		conf.Server.DebugPort = 0
		gin.SetMode(gin.ReleaseMode)
	}

	// create cacher
	cacher, err := cache.New(conf.Cacher)
	if err != nil {
		logrus.Fatalf("create cacher error: %s", err)
	}
	defer cacher.Close()

	digger, err := dig.New(conf.Digger, cacher)
	if err != nil {
		logrus.Fatal(err)
	}

	r := route.New(digger, conf.Server)
	r.Serve()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)

	<-sigChan
	logrus.Info("quit")

	return
}
