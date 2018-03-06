package dig

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/config"
	"github.com/wrfly/web-dns/dig/cache"
	"github.com/wrfly/web-dns/lib"
)

const (
	domainMaxLen = 64
	typeMaxLen   = 10
)

type Digger struct {
	cacher    cache.Cacher
	nsserver  []string
	blacklist []string
	timeout   time.Duration
	hijack    bool
	hosts     map[string]lib.Answer
}

func New(conf config.DiggerConfig, cacher cache.Cacher) (Digger, error) {
	logrus.Info("create new digger")
	hijack := false
	hosts := make(map[string]lib.Answer, 0)
	if conf.HostsFile != "" {
		if f, err := os.Open(conf.HostsFile); err == nil {
			logrus.Infof("use hijack file: %s", conf.HostsFile)
			hijack = true
			hosts = hostsHandler(f)
		}
	}
	return Digger{
		cacher:   cacher,
		nsserver: conf.DNS,
		timeout:  conf.Timeout,
		hijack:   hijack,
		hosts:    hosts,
	}, nil
}

func (d Digger) Dig(ctext context.Context, domain, typ string) ([]string, error) {
	answer := d.DigJson(ctext, domain, typ)
	return answer.IPs(), answer.Err
}

func hostKey(domain, typ string) string {
	return fmt.Sprintf("%s%s", domain, typ)
}

func (d Digger) DigJson(ctext context.Context, domain, typ string) (ans lib.Answer) {
	logrus.Debugf("digger: %s %s", domain, typ)

	// validate
	if len(domain) > domainMaxLen || len(typ) > typeMaxLen {
		return lib.Answer{
			Err: fmt.Errorf("u r kidding me, domain or type toooo long"),
		}
	}

	// hijack
	if d.hijack {
		if ans, exist := d.hosts[hostKey(domain, typ)]; exist {
			return ans
		}
	}

	var err error
	if ans, err = d.cacher.Get(domain, typ); err == nil {
		logrus.Debugf("return answer: %v", ans)
		return ans
	} else {
		logrus.Debugf("cacher error: %s", err)
	}

	first := make(chan lib.Answer)
	defer close(first)

	ctx, cancel := context.WithCancel(ctext)
	defer cancel()

	for _, ns := range d.nsserver {
		go func(ns string) {
			defer func() {
				x := recover()
				if x != nil {
					logrus.Errorf("got panic: %s", x)
				}
			}()
			// question always return
			r := lib.Question(ns, domain, typ, d.timeout)
			logrus.Debugf("%s got ip of %s: %v", ns, domain, r.IPs())
			if ctx.Err() != nil {
				logrus.Debug("abort answer")
				return
			}
			first <- r
		}(ns)
	}

	select {
	case ans = <-first:
		cancel()
		logrus.Infof("got answer: %v", ans)
		if err := d.cacher.Set(domain, typ, ans); err != nil {
			logrus.Errorf("set cache error: %s", err)
		}
	case <-ctx.Done():
		ans = lib.Answer{Err: fmt.Errorf("cancled")}
	}

	return ans
}

func simpleAnswer(ip string) lib.Answer {
	return lib.Answer{
		Result: []lib.Resp{
			{
				IP:  ip,
				TTL: 233,
			},
		},
	}
}

func hostsHandler(r io.ReadCloser) map[string]lib.Answer {
	var (
		lnum   = 0
		lines  = bufio.NewReader(r)
		hosts  = make(map[string]lib.Answer)
		domain string
		typ    string
		ip     string
	)
	for {
		lnum++
		s, err := lines.ReadString('\n')
		if err != nil && s == "" {
			if err == io.EOF {
				break
			}
		}

		if s[0] == '#' {
			continue
		}

		host := strings.Split(s, " ")
		hostValid := true
		switch len(host) {
		case 2:
			domain = host[0]
			typ = "A"
			ip = host[1]
		case 3:
			domain = host[0]
			typ = host[1]
			ip = host[2]
		default:
			hostValid = false
		}
		if hostValid {
			hosts[hostKey(domain, typ)] = simpleAnswer(ip)
		}
	}
	return hosts
}
