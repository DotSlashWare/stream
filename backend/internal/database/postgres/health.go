package postgres

import "time"

// Performs a health check on the PostgreSQL connection and returns the health status and latency.
func (m *Manager) GetHealth() (bool, int64) {
	start_time := time.Now()
	if !m.IsHealthy() {
		return false, 0
	}

	latency := time.Since(start_time).Milliseconds()
	return true, latency
}
