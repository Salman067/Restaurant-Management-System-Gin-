package middlewares

import (
	"fmt"
	"pi-inventory/common/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func LogAccessLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		protocol := ctx.Request.Proto
		xForwarded := ctx.Request.Header.Get("X-Forwarded-For")
		userAgent := ctx.Request.Header.Get("User-Agent")

		ctx.Next()

		end := time.Now()
		latency := strconv.FormatInt(end.Sub(start).Milliseconds(), 10)
		ip := ctx.ClientIP()
		bytesReceived := strconv.FormatInt(ctx.Request.ContentLength, 10)
		bytesSent := strconv.Itoa(ctx.Writer.Size())
		statusCode := strconv.Itoa(ctx.Writer.Status())

		formattedStart := start.Format(time.RFC3339)
		accessLogString := fmt.Sprintf("[%s] \"%s %s %s\" %s 0 %s %s %s \"%s\" %s \"%s\"", formattedStart, method, path,
			protocol, statusCode, bytesReceived, bytesSent, latency, xForwarded, userAgent, ip)
		logger.LogInfo(accessLogString)
	}
}
