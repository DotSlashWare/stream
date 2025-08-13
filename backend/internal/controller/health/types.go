package health

// Represents the basic health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// DetailedHealthResponse represents the detailed health check response
type DetailedHealthResponse struct {
	Status    string    `json:"status"`
	Timestamp string    `json:"timestamp"`
	Uptime    string    `json:"uptime"`
	Databases Databases `json:"databases"` // Contains health status of individual databases
}

// Represents the health status of multiple services
type Databases struct {
	Postgres ServiceHealth `json:"postgres"`
}

// Represents the health status of an individual service
type ServiceHealth struct {
	Status  string `json:"status"`
	Latency int64  `json:"latency_ms"`
	Uptime  string `json:"uptime,omitempty"` // Optional field for uptime in string
}
