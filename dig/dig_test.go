package dig

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/dig/cache"
)

func TestDig(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	ns := []string{
		"114.114.114.114:53",
		"8.8.8.8:53",
	}
	domain := "kfd.me"
	c, _ := cache.New("mem", "")
	digger, err := New(ns, time.Second, c)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	t.Run("Type A of kfd.me", func(t *testing.T) {
		ips, err := digger.Dig(ctx, domain, "A")
		if err != nil {
			t.Error(err)
		}
		for _, ip := range ips {
			t.Log(ip)
		}
	})
	t.Run("Type MX of kfd.me", func(t *testing.T) {
		ips, err := digger.Dig(ctx, domain, "MX")
		if err != nil {
			t.Error(err)
		}
		for _, ip := range ips {
			t.Log(ip)
		}
	})
	t.Run("Type X of kfd.me", func(t *testing.T) {
		_, err := digger.Dig(ctx, domain, "X")
		if err.Error() != "type not support" {
			t.Error(err)
		}
	})

	t.Run("Json Type A of kfd.me", func(t *testing.T) {
		r := digger.DigJson(ctx, domain, "A")
		if r.Err != nil {
			t.Error(r.Err)
		}
		t.Log(r)
	})
	t.Run("Json Type MX of kfd.me", func(t *testing.T) {
		r := digger.DigJson(ctx, domain, "A")
		if r.Err != nil {
			t.Error(r.Err)
		}
		t.Log(r)
	})
}
