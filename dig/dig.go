package dig

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	_ "github.com/wrfly/web-dns/dig/cache"
	"github.com/wrfly/web-dns/lib"
)

type Cacher interface {
	Set(domain string, ttl int) error
	Get(domain string) error
}

type Digger struct {
	cacher    Cacher
	nsserver  []string
	blacklist []string
	timeout   time.Duration
}

func New(cacher string, nsserver []string, timeout time.Duration) Digger {
	logrus.Info("create new digger")
	return Digger{
		nsserver: nsserver,
		timeout:  timeout,
	}
}

func (d Digger) Dig(ctext context.Context, domain, typ string) ([]string, error) {
	answer := d.DigJson(ctext, domain, typ)
	return answer.IPs(), answer.Err
}

func (d Digger) DigJson(ctext context.Context, domain, typ string) lib.Answer {
	logrus.Debugf("digger: %s %s", domain, typ)
	first := make(chan lib.Answer, 1)
	defer close(first)

	ctx, cancel := context.WithTimeout(ctext, d.timeout)
	defer cancel()

	for _, ns := range d.nsserver {
		go func(ns string) {
			defer func() {
				x := recover()
				if x != nil {
					logrus.Errorf("got panic: %s", x)
				}
			}()
			r := lib.Question(ns, domain, typ)
			logrus.Debugf("%s got ip of %s: %v", ns, domain, r.IPs())
			if ctx.Err() != nil {
				logrus.Debug("abort answer")
				return
			}
			first <- r
		}(ns)
	}
	return <-first
}
