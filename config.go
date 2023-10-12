package tlscertreloader

import "time"

type config struct {
	period time.Duration
}

type Option func(c *config)

func WithPeriod(t time.Duration) Option {
	return func(c *config) {
		c.period = t
	}
}
