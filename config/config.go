package config

import (
	"fmt"
)

type Config struct {
	Port      int
	DNS       []string
	CacheType string
	Limit     bool
	Rate      int
}

func (c Config) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("Port must be positive: %d", c.Port)
	}
	if len(c.DNS) == 0 {
		return fmt.Errorf("DNS upstream can not be empty")
	}
	return nil
}
