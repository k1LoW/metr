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
	metrics []Metric
}

// NewMetrics returns *Metrics
func NewMetrics() *Metrics {
	m := &Metrics{
		metrics: availableMetrics(),
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

// List returns metric info list
func (m *Metrics) List() []Metric {
	return m.metrics
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
func (m *Metrics) Each(f func(key string, value interface{}, format string)) {
	for _, metric := range m.metrics {
		if value, ok := m.Load(metric.Name); ok {
			f(metric.Name, value, metric.Format)
		} else {
			f(metric.Name, 0, metric.Format)
		}
	}
}

// Get returns metrics
func Get(i time.Duration) (*Metrics, error) {
	return getMetrics(i)
}
