package cache

import (
	"context"

	"github.com/wrfly/web-dns/lib"
)

type boltDBCacher struct{}

func (c *boltDBCacher) Set(ctx context.Context, ans lib.Answer) error {
	return nil
}

func (c *boltDBCacher) Get(ctx context.Context, domain, typ string) (lib.Answer, error) {
	return lib.Answer{}, nil
}
