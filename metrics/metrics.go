package metrics

import (
	"errors"
	"reflect"
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
	procPIDs        []int32
	procMetrics     []Metric
}

// NewMetrics returns *Metrics
func NewMetrics(interval time.Duration) *Metrics {
	m := &Metrics{
		collectInterval: interval,
		metrics:         AvailableMetrics(),
	}
	return m
}

func (m *Metrics) SetProcPIDs(pids []int32) error {
	for _, pid := range pids {
		if pid <= 0 {
			return errors.New("PID should be >= 0")
		}
		m.procPIDs = append(m.procPIDs, pid)
	}
	if len(pids) > 0 {
		m.procMetrics = AvailableProcMetrics()
	}

	return nil
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
	// process metrics
	for _, metric := range m.procMetrics {
		if value, ok := m.Load(metric.Name); ok {
			metrics[metric.Name] = value
		} else {
			metrics[metric.Name] = 0
		}
	}
	// system metrics
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
	// process metrics
	for _, metric := range m.procMetrics {
		if value, ok := m.Load(metric.Name); ok {
			f(metric, value)
		} else {
			f(metric, 0)
		}
	}
	// system metrics
	for _, metric := range m.metrics {
		if value, ok := m.Load(metric.Name); ok {
			f(metric, value)
		} else {
			f(metric, 0)
		}
	}
}

// GetMetrics returns metrics
func GetMetrics(interval time.Duration, pids ...int32) (*Metrics, error) {
	m := NewMetrics(interval)
	if len(pids) > 0 && !reflect.DeepEqual(pids, []int32{0}) {
		err := m.SetProcPIDs(pids)
		if err != nil {
			return nil, err
		}
	}
	err := m.Collect()
	if err != nil {
		return nil, err
	}
	return m, nil
}
