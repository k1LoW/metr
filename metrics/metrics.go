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
	interval time.Duration
	metrics  []Metric
}

// NewMetrics returns *Metrics
func NewMetrics(i time.Duration) *Metrics {
	m := &Metrics{
		interval: i,
		metrics:  AvailableMetrics(),
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

// Collect returns metrics
func Collect(i time.Duration) (*Metrics, error) {
	return NewMetrics(i).Collect()
}
