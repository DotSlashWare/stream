package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// Custom Formatter struct for GIN log output. It represents the gin.LogFormatter interface.
type CustomFormatter struct{}

// Formats the GIN log output according to a custom format.
func (formatter *CustomFormatter) Format(param gin.LogFormatterParams) string {
	timestamp := param.TimeStamp.Format("2006/01/02 - 15:04:05")

	// @returns: "[2006/01/02 - 15:04:05] - [status code] | [latency] | [client IP] | [method] [path]"
	return fmt.Sprintf("[%s] - %d | %v | %s | %s %s\n",
		timestamp,
		param.StatusCode,
		param.Latency,
		param.ClientIP,
		param.Method,
		param.Path,
	)
}
