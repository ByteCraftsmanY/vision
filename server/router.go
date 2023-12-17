package server

import (
	"github.com/gin-gonic/gin"
	"vision/controllers"
	"vision/middlewares"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	health := new(controllers.HealthController)
	router.GET("/health", health.Status)

	//router.Use(middlewares.AuthMiddleware())
	//router.Use(middlewares.LogMiddleware())

	v1 := router.Group("v1")
	{
		authGroup := v1.Group("auth")
		{
			auth := new(controllers.AuthController)
			authGroup.POST("/", auth.Login)
		}

		userGroup := v1.Group("user")
		{
			user := new(controllers.UserController)
			//userGroup.Get("/:id", user.Retrieve)
			userGroup.GET("/", user.Retrieve, middlewares.AuthMiddleware())
			userGroup.POST("/", user.Store)
		}
	}
	return router
}
