package cache

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/web-dns/lib"
)

// cacher type
const (
	MemCache   = "mem"
	RedisCache = "redis"
	BoltCache  = "bolt"
)

type Cacher interface {
	Set(domain, typ string, ans lib.Answer) error
	Get(domain, typ string) (lib.Answer, error)
}

func cacheKey(domain, typ string) string {
	return fmt.Sprintf("%s-%s", domain, typ)
}
func New(cacheTyp string, addr ...string) (Cacher, error) {
	logrus.Debugf("new cacher: %s", cacheTyp)
	switch cacheTyp {
	case MemCache:
		return &memCacher{
			n:       MemCache,
			storage: make(memKVStorage),
		}, nil
	case RedisCache:
		client := redis.NewClient(&redis.Options{
			Addr:     addr[0],
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		if _, err := client.Ping().Result(); err != nil {
			return nil, err
		}
		return &redisCacher{cli: client}, nil
	default:
		return nil, fmt.Errorf("cache type [%s] not support", cacheTyp)
	}
}
