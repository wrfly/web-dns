package cache

import (
	"context"
	"fmt"

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
	Set(ctx context.Context, domain, typ string, ans lib.Answer) error
	Get(ctx context.Context, domain, typ string) (lib.Answer, error)
}

func New(cacheTyp string) (Cacher, error) {
	logrus.Debugf("new cacher: %s", cacheTyp)
	switch cacheTyp {
	case MemCache:
		return &memCacher{
			n:       MemCache,
			storage: make(memKVStorage),
		}, nil
	default:
		return nil, fmt.Errorf("cache type [%s] not support", cacheTyp)
	}
}
