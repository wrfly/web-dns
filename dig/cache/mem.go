package cache

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/lib"
)

type memCacher struct {
	n       string
	m       sync.RWMutex
	storage memKVStorage
}

type memKVStorage map[string]lib.Answer

func (c *memCacher) Set(domain, typ string, ans lib.Answer) error {
	key := cacheKey(domain, typ)
	c.m.Lock()
	defer c.m.Unlock()
	c.storage[key] = ans
	logrus.Debugf("cacher %s set %s=%v", c.n, key, ans)
	return nil
}

func (c *memCacher) Get(domain, typ string) (lib.Answer, error) {
	c.m.RLock()
	defer c.m.RUnlock()
	key := cacheKey(domain, typ)
	if ansGot, got := c.storage[key]; got {
		// TODO: ugly here
		ans := lib.Answer{
			Result: make([]lib.Resp, len(ansGot.Result)),
		}
		ans.DigAt = ansGot.DigAt
		copy(ans.Result, ansGot.Result)
		logrus.Debugf("cacher %s get %s=%v", c.n, key, ans)
		x := uint32(time.Now().Unix() - ans.DigAt)
		logrus.Debugf("x=%d", x)
		for i := range ans.Result {
			if ans.Result[i].TTL >= x {
				ans.Result[i].TTL -= x
			} else {
				return ans, fmt.Errorf("answer out of date")
			}
		}
		return ans, nil
	}
	return lib.Answer{}, fmt.Errorf("404")
}

func (c *memCacher) Close() error {
	return nil
}
