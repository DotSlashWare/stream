package movie

import (
	"time"

	"github.com/artumont/DotSlashStream/backend/internal/service/tmdb"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	initTime    time.Time
	tmdbService *tmdb.Service
}

func NewMovieController(tmdbService *tmdb.Service) *Controller {
	return &Controller{
		initTime:    time.Now(),
		tmdbService: tmdbService,
	}
}

func (c *Controller) Register(router *gin.Engine) {
	movieGroup := router.Group("/movies")
	{
		movieGroup.GET("", nil)
	}
}
