package dig

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/config"
	"github.com/wrfly/web-dns/dig/cache"
)

func TestDig(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	conf := config.DiggerConfig{
		DNS: []string{
			"1.1.1.1:53",
			"8.8.8.8:53",
		},
		Timeout: time.Second,
	}
	domain := "kfd.me"
	cacher, _ := cache.New(config.CacherConfig{CacheType: "mem"})
	digger, err := New(conf, cacher)
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()

	for _, typ := range []string{"A", "AAAA", "MX", "NS", "TXT"} {
		t.Logf("Type %s of kfd.me", typ)
		r := digger.DigJSON(ctx, domain, "A")
		if r.Err != nil {
			t.Error(r.Err)
		}
		t.Log(r)
	}

}
