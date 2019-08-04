package metrics

import (
	"sync"
	"time"
)

type Metrics struct {
	sync.Map
	keys []string
}

// NewMetrics ...
func NewMetrics() *Metrics {
	m := &Metrics{
		keys: []string{
			"cpu",
			"mem",
			"swap",
			"user",
			"system",
			"idle",
			"nice",
			"iowait",
			"irq",
			"softirq",
			"steal",
			"guest",
			"guest_nice",
			"stolen",
			"load1",
			"load5",
			"load15",
		},
	}
	return m
}

func (m *Metrics) Each(f func(key string, value interface{})) {
	for _, key := range m.keys {
		if value, ok := m.Load(key); ok {
			f(key, value)
		} else {
			f(key, 0)
		}
	}
}

// Get returns metrics
func Get(i time.Duration) (*Metrics, error) {
	return getMetrics(i)
}
