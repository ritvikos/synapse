package robots

import (
	"errors"
	"time"
)

// Configures the [RobotsResolver] instance
type RobotsConfig struct {
	// User-agent to be used when fetching robots.txt from hosts
	UserAgent string

	// Default TTL for cached robots.txt entries
	TTL time.Duration
}

func (c *RobotsConfig) validate() error {
	if c.TTL <= 0 {
		return ErrRobotsInvalidTTL
	}
	if c.UserAgent == "" {
		return errors.New("robots resolver: user-agent cannot be empty")
	}
	return nil
}
