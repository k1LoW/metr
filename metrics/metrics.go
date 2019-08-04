package metrics

import (
	"sync"
	"time"
)

type Metric struct {
	Name        string
	Description string
	Format      string
}

// Metrics struct
type Metrics struct {
	sync.Map
	metrics []Metric
}

// NewMetrics returns *Metrics
func NewMetrics() *Metrics {
	m := &Metrics{
		metrics: []Metric{
			{"cpu", "Percentage of cpu used.", "%f"},
			{"mem", "Percentage of RAM used.", "%f"},
			{"swap", "Amount of memory that has been swapped out to disk (bytes).", "%d"},

			{"user", "Percentage of CPU utilization that occurred while executing at the user level.", "%f"},
			{"system", "Percentage of CPU utilization that occurred while executing at the system level.", "%f"},
			{"idle", "Percentage of time that CPUs were idle and the system did not have an outstanding disk I/O request.", "%f"},
			{"nice", "Percentage of CPU utilization that occurred while executing at the user level with nice priority.", "%f"},
			{"iowait", "Percentage of time that CPUs were idle during which the system had an outstanding disk I/O request.", "%f"},
			{"irq", "Percentage of time spent by CPUs to service hardware interrupts.", "%f"},
			{"softirq", "Percentage of time spent by CPUs to service software interrupts.", "%f"},
			{"steal", "Percentage of time spent in involuntary wait by the virtual CPUs while the hypervisor was servicing another virtual processor.", "%f"},
			{"guest", "Percentage of time spent by CPUs to run a virtual processor.", "%f"},
			{"guest_nice", "Percentage of time spent by CPUs to run a virtual processor with nice priority.", "%f"},

			{"load1", "Load avarage for 1 minute.", "%f"},
			{"load5", "Load avarage for 5 minutes.", "%f"},
			{"load15", "Load avarage for 15 minutes.", "%f"},
		},
	}
	return m
}

// List returns metric info list
func (m *Metrics) List() []Metric {
	return m.metrics
}

// Each returns ordered metrics
func (m *Metrics) Each(f func(key string, value interface{})) {
	for _, metric := range m.metrics {
		if value, ok := m.Load(metric.Name); ok {
			f(metric.Name, value)
		} else {
			f(metric.Name, 0)
		}
	}
}

// Get returns metrics
func Get(i time.Duration) (*Metrics, error) {
	return getMetrics(i)
}
