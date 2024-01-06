package server

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
	"vision/controllers"
	"vision/middlewares"
)

func NewRouter(logger *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))

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

		cctvGroup := v1.Group("cctv")
		{
			cctv := new(controllers.CCTVController)
			cctvGroup.GET("/:ID", cctv.Get)
			cctvGroup.Any("/all", cctv.GetAll)
			cctvGroup.POST("/", cctv.Add)
			cctvGroup.PUT("/", cctv.Update)
			cctvGroup.DELETE("/:ID", cctv.Delete)
		}

		organizationGroup := v1.Group("organization")
		{
			organization := new(controllers.OrganizationController)
			organizationGroup.GET("/:ID", organization.Get)
			organizationGroup.Any("/all", organization.GetAll)
			organizationGroup.POST("/", organization.Add)
			organizationGroup.PUT("/", organization.Update)
			organizationGroup.DELETE("/:ID", organization.Delete)
		}

	}
	return router
}
