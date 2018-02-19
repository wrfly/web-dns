package dig

import (
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
}

func New(cacher string, nsserver []string) Digger {
	logrus.Info("create new digger")
	return Digger{
		nsserver: nsserver,
	}
}

func (d Digger) Dig(domain, typ string) ([]string, error) {
	logrus.Debugf("digger: %s %s", domain, typ)
	first := make(chan []string, 1)
	errChan := make(chan error, 1)
	defer close(first)
	defer close(errChan)

	for _, ns := range d.nsserver {
		go func(ns string) {
			defer func() {
				x := recover()
				if x != nil {
					logrus.Errorf("got panic: %s", x)
				}
			}()
			r := lib.Question(ns, domain, typ)
			if err := r.Err; err != nil {
				errChan <- err
				first <- nil
			}
			logrus.Debugf("%s got ip of %s: %v", ns, domain, r.IPs())
			errChan <- nil
			first <- r.IPs()
		}(ns)
	}
	return <-first, <-errChan
}

func (d Digger) DigJson(domain, typ string) lib.Answer {
	logrus.Debugf("digger: %s %s", domain, typ)
	first := make(chan lib.Answer, 1)
	errChan := make(chan error, 1)
	defer close(first)
	defer close(errChan)

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
			first <- r
		}(ns)
	}
	return <-first
}
