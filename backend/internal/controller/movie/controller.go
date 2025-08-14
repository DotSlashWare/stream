package movie

import (
	"time"

	"github.com/artumont/DotSlashStream/backend/internal/service/tmdb"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	initTime     time.Time
	tmdbService  *tmdb.Service
	streamApiUrl string
}

func NewMovieController(tmdbService *tmdb.Service, streamApiUrl string) *Controller {
	return &Controller{
		initTime:    time.Now(),
		tmdbService: tmdbService,
		streamApiUrl: streamApiUrl,
	}
}

func (controller *Controller) Register(router *gin.Engine) {
	movieGroup := router.Group("/movie")
	{
		movieGroup.GET("/search", controller.SearchForMovie)
		movieGroup.GET("/id/:id", controller.GetMovieById)
	}
}
