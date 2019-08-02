package metrics

// Get returns metrics
func Get() (map[string]interface{}, error) {
	return getMetrics()
}
