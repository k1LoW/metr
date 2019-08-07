package metrics

import (
	"sync"
	"time"
)

type Metric struct {
	Name        string
	Description string
	Format      string
	Unit        string
}

// Metrics struct
type Metrics struct {
	sync.Map
	collectInterval time.Duration
	metrics         []Metric
}

// NewMetrics returns *Metrics
func NewMetrics(i time.Duration) *Metrics {
	m := &Metrics{
		collectInterval: i,
		metrics:         AvailableMetrics(),
	}
	return m
}

func (m *Metrics) Format(key string) string {
	for _, metric := range m.metrics {
		if metric.Name == key {
			return metric.Format
		}
	}
	return "%v"
}

// Raw returns raw metrics map
func (m *Metrics) Raw() map[string]interface{} {
	metrics := map[string]interface{}{}
	for _, metric := range m.metrics {
		if value, ok := m.Load(metric.Name); ok {
			metrics[metric.Name] = value
		} else {
			metrics[metric.Name] = 0
		}
	}
	return metrics
}

// Each returns ordered metrics
func (m *Metrics) Each(f func(metric Metric, value interface{})) {
	for _, metric := range m.metrics {
		if value, ok := m.Load(metric.Name); ok {
			f(metric, value)
		} else {
			f(metric, 0)
		}
	}
}

// GetMetrics returns metrics
func GetMetrics(i time.Duration) (*Metrics, error) {
	m := NewMetrics(i)
	err := m.Collect()
	if err != nil {
		return nil, err
	}
	return m, nil
}
