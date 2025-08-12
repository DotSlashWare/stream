package bootstrap

import (
	//"log"
	"time"
	"github.com/gin-gonic/gin"
)

type Application struct {
	InitTime time.Time
	Env      *Env
	Router   *gin.Engine
}
