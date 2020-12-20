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

// Digger ...
type Digger struct {
	cacher    cache.Cacher
	nsServer  []string
	blacklist []string
	timeout   time.Duration
	hijack    bool
	hosts     map[string]lib.Answer
}

// New DNS digger
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
		nsServer: conf.DNS,
		timeout:  conf.Timeout,
		hijack:   hijack,
		hosts:    hosts,
	}, nil
}

// Dig ...
func (d Digger) Dig(ctx context.Context, domain, typ string) ([]string, error) {
	answer := d.DigJSON(ctx, domain, typ) // filter DNS type
	return answer.Hosts(), answer.Err
}

func hostKey(domain, typ string) string {
	return fmt.Sprintf("%s%s", domain, typ)
}

// DigJSON returns json
func (d Digger) DigJSON(ctx context.Context, domain, typ string) lib.Answer {
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

	if ans, err := d.cacher.Get(domain, typ); err != nil {
		logrus.Debugf("cacher error: %s", err)
	} else {
		logrus.Debugf("return answer: %v", ans)
		return ans
	}

	first := make(chan lib.Answer)
	defer close(first)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, ns := range d.nsServer {
		go func(ns string) {
			// question always return
			r := lib.Question(lib.QueryOption{
				NSServer: ns,
				Domain:   domain,
				Timeout:  d.timeout,
				Type:     typ,
			})
			logrus.Debugf("%s got ip of %s: %v", ns, domain, r.Hosts())
			if ctx.Err() != nil {
				logrus.Debug("abort answer")
				return
			}
			select {
			case first <- r:
			default:
			}
		}(ns)
	}

	tm := time.NewTimer(d.timeout)
	defer tm.Stop()

	select {
	case <-tm.C:
		return lib.Answer{Err: fmt.Errorf("timeout")}

	case ans := <-first:
		cancel()
		logrus.Debugf("got answer: %v", ans)
		if err := d.cacher.Set(domain, typ, ans); err != nil {
			logrus.Errorf("set cache error: %s", err)
		}
		return ans

	case <-ctx.Done():
		return lib.Answer{Err: fmt.Errorf("cancled")}
	}
}

func simpleAnswer(ip, typ string) lib.Answer {
	return lib.Answer{
		Result: []lib.Result{
			{
				Host: ip,
				Type: typ,
				TTL:  233,
			},
		},
	}
}

func hostsHandler(r io.ReadCloser) map[string]lib.Answer {
	var (
		lNum   = 0
		lines  = bufio.NewReader(r)
		hosts  = make(map[string]lib.Answer)
		domain string
		typ    string
		ip     string
	)
	for {
		lNum++
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
			hosts[hostKey(domain, typ)] = simpleAnswer(ip, typ)
		}
	}
	return hosts
}
