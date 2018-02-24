package cache

import (
	"context"

	"github.com/wrfly/web-dns/lib"
)

type redisCacher struct{}

func (c *redisCacher) Set(ctx context.Context, ans lib.Answer) error {
	return nil
}

func (c *redisCacher) Get(ctx context.Context, domain, typ string) (lib.Answer, error) {
	return lib.Answer{}, nil
}
