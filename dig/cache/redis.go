package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	"github.com/wrfly/web-dns/lib"
)

type redisCacher struct {
	cli *redis.Client
}

func (c *redisCacher) Set(domain, typ string, ans lib.Answer) error {
	key := cacheKey(domain, typ)
	var ttl time.Duration
	for _, r := range ans.Result {
		ttl = time.Second * time.Duration(r.TTL)
		break
	}
	logrus.Debugf("redis set: %s=%v ttl: %v", key, ans, ttl)
	s := c.cli.Set(key, ans.Marshal(), 0)
	return s.Err()
}

func (c *redisCacher) Get(domain, typ string) (lib.Answer, error) {
	ans := lib.Answer{}
	s := c.cli.Get(cacheKey(domain, typ))
	bs, err := s.Bytes()
	if err != nil {
		return ans, err
	}
	json.Unmarshal(bs, &ans)
	logrus.Debugf("redis get: %v", ans)
	return ans, nil
}

func (c *redisCacher) Close() error {
	return c.cli.Close()
}
