package config

import (
	"time"
)

type Config struct {
	Port      int
	DNS       []string
	BLK       []string // black list
	CacheType string
	Limit     bool
	Rate      int
	Timeout   time.Duration
}
