package config

type Config struct {
	Port      int
	DNS       []string
	CacheType string
	Limit     bool
	Rate      int
}
