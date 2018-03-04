package cache

import (
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"

	"github.com/wrfly/web-dns/config"
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
	Close() error
}

func cacheKey(domain, typ string) string {
	return fmt.Sprintf("%s-%s", domain, typ)
}
func New(conf config.CacherConfig) (Cacher, error) {
	logrus.Debugf("new cacher: %s", conf.CacheType)
	switch conf.CacheType {
	case MemCache:
		return &memCacher{
			n:       MemCache,
			storage: make(memKVStorage),
		}, nil
	case RedisCache:
		client := redis.NewClient(&redis.Options{
			Addr:     conf.RedisAddr,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		if _, err := client.Ping().Result(); err != nil {
			return nil, err
		}
		return &redisCacher{cli: client}, nil
	case BoltCache:
		dbPath := "dns.db"
		os.Remove(dbPath)
		db, err := bolt.Open(dbPath, 0600, nil)
		if err != nil {
			return nil, err
		}
		if err := db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucket([]byte("dns"))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
			return nil
		}); err != nil {
			return nil, err
		}
		return &boltDBCacher{dbPath: dbPath, db: db, bktName: []byte("dns")}, nil
	default:
		return nil, fmt.Errorf("cache type [%s] not support", conf.CacheType)
	}
}
