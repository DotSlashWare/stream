package profile

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	initTime     time.Time
}

func NewMovieController() *Controller {
	return &Controller{
		initTime:    time.Now(),
	}
}

func (controller *Controller) Register(router *gin.Engine) {
	/*
	profileGroup := router.Group("/profile")
	{
		profileGroup.GET("/", nil)
	}
	*/
}
