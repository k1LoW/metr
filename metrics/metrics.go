package metrics

import "time"

// Get returns metrics
func Get(i time.Duration) (map[string]interface{}, error) {
	return getMetrics(i)
}
