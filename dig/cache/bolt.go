package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"

	"github.com/wrfly/web-dns/lib"
)

type boltDBCacher struct {
	dbPath  string
	db      *bolt.DB
	bktName []byte
}

func (c *boltDBCacher) Set(domain, typ string, ans lib.Answer) error {
	return c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bktName)
		key := cacheKey(domain, typ)
		logrus.Debugf("cacher bolt set %s=%v", key, ans)
		err := b.Put([]byte(key), ans.Marshal())
		return err
	})
}

func (c *boltDBCacher) Get(domain, typ string) (ans lib.Answer, err error) {
	err = c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(c.bktName)
		key := cacheKey(domain, typ)
		v := b.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("not found")
		}
		ansGot := lib.Answer{}
		if err := json.Unmarshal(v, &ansGot); err != nil {
			return err
		}
		// TODO: ugly here
		ans = lib.Answer{
			Result: make([]lib.Resp, len(ansGot.Result)),
		}
		ans.DigAt = ansGot.DigAt
		copy(ans.Result, ansGot.Result)
		logrus.Debugf("cacher bolt get %s=%v", key, ans)
		x := uint32(time.Now().Unix() - ans.DigAt)
		for i := range ans.Result {
			if ans.Result[i].TTL >= x {
				ans.Result[i].TTL -= x
			} else {
				return fmt.Errorf("answer out of date")
			}
		}
		return nil
	})
	return
}

func (c *boltDBCacher) Close() error {
	if err := c.db.Close(); err != nil {
		return err
	}
	return os.Remove(c.dbPath)
}
