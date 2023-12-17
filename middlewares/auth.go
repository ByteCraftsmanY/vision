package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vision/utils/token"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := token.Validate(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
