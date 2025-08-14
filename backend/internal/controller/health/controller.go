package health

import (
	"time"

	"github.com/artumont/DotSlashStream/backend/internal/database/postgres"
	"github.com/gin-gonic/gin"
)

// Controller for managing health checks of various database connections and the app itself.
type Controller struct {
	initTime        time.Time
	postgresManager *postgres.Manager
}

// Factory function to create a new Controller instance.
func NewHealthController(
	initTime time.Time,
	postgresManager *postgres.Manager,
) *Controller {
	return &Controller{
		initTime:        initTime,
		postgresManager: postgresManager,
	}
}

// Sets up the routes for the health controller.
func (controller *Controller) Register(router *gin.Engine) {
	healthGroup := router.Group("/health")
	{
		healthGroup.GET("", controller.GetHealth)
		healthGroup.GET("/detailed", controller.GetHealthDetailed)
	}
}
