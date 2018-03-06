package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"

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
	app := cli.App{
		Name:    "web-dns",
		Usage:   "Query domain via HTTP(S)",
		Version: "0.1",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "wrfly",
				Email: "mr.wrfly@gmail.com",
			},
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Usage:       "port to listen on",
				EnvVars:     env_vars("port"),
				Value:       8080,
				Destination: &conf.Server.Port,
			},
			&cli.StringSliceFlag{
				Name:    "dns",
				Usage:   "dns servers",
				EnvVars: env_vars("dns"),
				Value:   cli.NewStringSlice("8.8.8.8:53", "8.8.4.4:53"),
			},
			&cli.IntFlag{
				Name:    "timeout",
				Usage:   "dig timeout (second)",
				EnvVars: env_vars("timeout"),
				Value:   1,
			},
			&cli.IntFlag{
				Name:        "rate",
				Usage:       "rate of requests per minute per IP",
				EnvVars:     env_vars("rate"),
				Value:       1000,
				Destination: &conf.Server.Rate,
			},
			&cli.StringFlag{
				Name:        "cache",
				Usage:       "cache type: mem|redis|bolt",
				Value:       "mem",
				EnvVars:     env_vars("cache"),
				Destination: &conf.Cacher.CacheType,
			},
			&cli.StringFlag{
				Name:        "redis-addr",
				Usage:       "this flag is used for redis cacher",
				Value:       "localhost:6379",
				EnvVars:     env_vars("redis-addr"),
				Destination: &conf.Cacher.RedisAddr,
			},
			&cli.StringSliceFlag{
				Name:    "black-list",
				Usage:   "blacklist of clients",
				EnvVars: env_vars("black-list"),
				Value:   cli.NewStringSlice("8.8.8.8", "4.4.4.4"),
			},
			&cli.StringFlag{
				Name:        "hosts",
				Usage:       "hijack hosts file path",
				Value:       "",
				EnvVars:     env_vars("hosts"),
				Destination: &conf.Digger.HostsFile,
			},
			&cli.BoolFlag{
				Name:        "debug",
				EnvVars:     env_vars("debug"),
				Usage:       "debug log-level, metrics and pprof debug",
				Destination: &conf.Debug,
			},
			&cli.IntFlag{
				Name:        "debug-port",
				EnvVars:     env_vars("debug-port"),
				Usage:       "server debug port",
				Value:       8081,
				Destination: &conf.Server.DebugPort,
			},
		},
		Action: func(c *cli.Context) error {
			tot := time.Second * time.Duration(c.Int("timeout"))
			conf.Digger.Timeout = tot
			conf.Server.BLK = c.StringSlice("black-list")
			conf.Digger.DNS = c.StringSlice("dns")
			if err := run(*conf); err != nil {
				logrus.Error(err)
			}
			return nil
		},
	}
	app.Run(os.Args)
}

func env_vars(n string) []string {
	return []string{
		// web dns config
		fmt.Sprintf("WDC_%s",
			strings.Replace(strings.ToUpper(n), "-", "_", -1)),
	}
}

func run(conf config.Config) error {
	if conf.Debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode")
	} else {
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

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)

	select {
	case <-sigChan:
		logrus.Info("quit")
		return nil
	case err := <-r.Serve():
		return err
	}
}
