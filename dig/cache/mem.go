package cache

import (
	"context"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/lib"
)

type memCacher struct {
	n       string
	m       sync.RWMutex
	storage memKVStorage
}

type memKVStorage map[string]lib.Answer

func (c *memCacher) Set(ctx context.Context,
	domain, typ string, ans lib.Answer) error {
	c.m.Lock()
	defer c.m.Unlock()
	key := fmt.Sprintf("%s-%s", domain, typ)
	c.storage[key] = ans
	logrus.Debugf("cacher %s set %s=%v", c.n, key, ans)
	return nil
}

func (c *memCacher) Get(ctx context.Context,
	domain, typ string) (lib.Answer, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	key := fmt.Sprintf("%s-%s", domain, typ)
	if ans, got := c.storage[key]; got {
		logrus.Debugf("cacher %s get %s=%v", c.n, key, ans)
		return ans, nil
	}
	return lib.Answer{}, fmt.Errorf("404")
}
