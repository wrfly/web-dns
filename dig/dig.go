package dig

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/dig/cache"
	"github.com/wrfly/web-dns/lib"
)

type Digger struct {
	cacher    cache.Cacher
	nsserver  []string
	blacklist []string
	timeout   time.Duration
}

func New(nsserver []string, timeout time.Duration, cacher cache.Cacher) (Digger, error) {
	logrus.Info("create new digger")
	return Digger{
		cacher:   cacher,
		nsserver: nsserver,
		timeout:  timeout,
	}, nil
}

func (d Digger) Dig(ctext context.Context, domain, typ string) ([]string, error) {
	answer := d.DigJson(ctext, domain, typ)
	return answer.IPs(), answer.Err
}

func (d Digger) DigJson(ctext context.Context, domain, typ string) (ans lib.Answer) {
	logrus.Debugf("digger: %s %s", domain, typ)
	var err error
	if ans, err = d.cacher.Get(domain, typ); err == nil {
		logrus.Debugf("return answer: %v", ans)
		return ans
	}

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

	ans = <-first
	if err := d.cacher.Set(domain, typ, ans); err != nil {
		logrus.Errorf("set cache error: %s", err)
	}

	return ans
}
