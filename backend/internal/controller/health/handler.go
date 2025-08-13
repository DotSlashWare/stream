package health

import (
	"time"

	"github.com/gin-gonic/gin"
)

func (controller *Controller) GetHealth(ctx *gin.Context) {
	// @returns
	// {
	//   "status": "healthy",
	//   "timestamp": "2023-10-01T12:00:00Z"
	// }
	ctx.JSON(200, HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func (controller *Controller) GetHealthDetailed(ctx *gin.Context) {
	// @returns
	// {
	//   "status": "healthy",
	//   "timestamp": "2023-10-01T12:00:00Z",
	//   "uptime": "1h30m",
	//   "databases": {
	//     "postgres": { "status": "healthy/unhealthy", "latency_ms": 20 },
	//   }
	// }

	postgresStatus, postgresLatency := controller.postgresManager.GetHealth()

	databases := Databases{
		Postgres: ServiceHealth{
			Status:  getHealthStatus(postgresStatus),
			Latency: postgresLatency,
			Uptime:  time.Since(controller.postgresManager.ConnectionTime).String(),
		},
	}

	ctx.JSON(200, DetailedHealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Uptime:    time.Since(controller.initTime).String(),
		Databases: databases,
	})
}

func getHealthStatus(isHealthy bool) string {
	if isHealthy {
		return "healthy"
	}
	return "unhealthy"
}
