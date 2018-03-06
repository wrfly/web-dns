package config

import (
	"time"
)

type Config struct {
	Digger DiggerConfig
	Server ServerConfig
	Cacher CacherConfig
	Debug  bool
}

type CacherConfig struct {
	CacheType string
	RedisAddr string
}

type ServerConfig struct {
	Port      int
	DebugPort int
	BLK       []string // black list
	Limit     bool
	Rate      int
}

type DiggerConfig struct {
	DNS       []string
	HostsFile string // for hijack
	Timeout   time.Duration
}

func New() *Config {
	return &Config{
		Digger: DiggerConfig{
			DNS: make([]string, 0),
		},
		Server: ServerConfig{
			BLK: make([]string, 0),
		},
	}
}
