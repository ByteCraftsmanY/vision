package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func LogMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()

		latencyTime := endTime.Sub(startTime)
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()

		zap.L().Info("HTTP REQUEST",
			zap.Int("STATUS", statusCode),
			zap.String("METHOD", reqMethod),
			zap.String("CLIENT_IP", clientIP),
			zap.String("REQUEST_URI", reqUri),
			zap.Duration("LATENCY", latencyTime),
		)

		ctx.Next()
	}
}
