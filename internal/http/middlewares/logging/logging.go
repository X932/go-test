package logging_middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		log.Printf("Logging ==== start")
		ctx.Next()

		durationTime := time.Since(startTime)
		httpStatusCode := ctx.Writer.Status()

		log.Printf("Logging ==== [%s] %s - %d - %v", method, path, httpStatusCode, durationTime)
	}
}
