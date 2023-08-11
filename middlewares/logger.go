package middlewares

import (
	"pi-inventory/common/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetTraceID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := ctx.Request.Header.Get("X-Trace-ID")
		if traceId == "" {
			traceId = uuid.New().String()
		}
		logger.WithField("traceid", traceId)
		ctx.Writer.Header().Add("X-Trace-ID", traceId)
		ctx.Next()
	}
}
