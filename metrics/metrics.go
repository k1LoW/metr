package metrics

import (
	"errors"
	"sync"
	"time"

	"github.com/shirou/gopsutil/process"
)

type Metric struct {
	Name         string
	Description  string
	Format       string
	Unit         string
	InitialValue interface{}
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
	}
	m.InitializeMetrics()
	return m
}

func (m *Metrics) InitializeMetrics() {
	m.metrics = AvailableMetrics()
	for _, metric := range m.metrics {
		m.Store(metric.Name, metric.InitialValue)
	}
}

func (m *Metrics) InitializeProcMetrics() {
	m.procMetrics = AvailableProcMetrics()
	for _, metric := range m.procMetrics {
		m.Store(metric.Name, metric.InitialValue)
	}
}

func (m *Metrics) SetProcPIDs(pids []int32) error {
	for _, pid := range pids {
		if pid <= 0 {
			return errors.New("PID should be >= 0")
		}
		m.procPIDs = append(m.procPIDs, pid)
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
	if len(pids) > 0 && !(len(pids) == 1 && pids[0] == 0) {
		return GetMetricsByPIDs(interval, pids)
	}
	m := NewMetrics(interval)
	err := m.Collect()
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetMetricsByPIDs returns metrics
func GetMetricsByPIDs(interval time.Duration, pids []int32) (*Metrics, error) {
	if len(pids) == 0 {
		return nil, errors.New("empty pids")
	}
	m := NewMetrics(interval)
	m.InitializeProcMetrics()
	err := m.SetProcPIDs(pids)
	if err != nil {
		return nil, err
	}
	err = m.Collect()
	if err != nil {
		return nil, err
	}
	return m, nil
}

// GetMetricsByName returns metrics
func GetMetricsByName(interval time.Duration, name string) (*Metrics, error) {
	if name == "" {
		return nil, errors.New("empty name")
	}
	m := NewMetrics(interval)
	m.InitializeProcMetrics()

	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	pids := []int32{}
	for _, p := range processes {
		pname, err := p.Name()
		if err != nil {
			continue
		}
		if pname == name {
			pids = append(pids, p.Pid)
		}
	}
	err = m.SetProcPIDs(pids)
	if err != nil {
		return nil, err
	}

	err = m.Collect()
	if err != nil {
		return nil, err
	}
	return m, nil
}
