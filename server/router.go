package server

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
	"vision/handlers"
	"vision/middlewares"
)

func NewRouter(params *Params) *gin.Engine {
	router := gin.New()
	router.Use(ginzap.RecoveryWithZap(params.Logger, true))
	router.Use(ginzap.Ginzap(params.Logger, time.RFC3339, true))

	health := handlers.NewHealthHandler()
	router.GET("/health", health.Status)

	//router.Use(middlewares.AuthMiddleware())
	//router.Use(middlewares.LogMiddleware())

	v1 := router.Group("v1")
	{
		authGroup := v1.Group("auth")
		{
			auth := handlers.NewAuthHandler()
			authGroup.POST("/", auth.Login)
			authGroup.POST("/initiate", auth.GenerateOTP)
			authGroup.POST("/verify", auth.VerifyOTP)
		}

		userGroup := v1.Group("user")
		{
			user := handlers.NewUserHandler()
			//userGroup.Get("/:id", user.Retrieve)
			userGroup.GET("/", user.Retrieve, middlewares.AuthMiddleware())
			userGroup.POST("/", user.Store)
		}

		cctvGroup := v1.Group("cctv")
		{
			cctv := handlers.NewProductHandler()
			cctvGroup.GET("/:ID", cctv.Get)
			//cctvGroup.Any("/all", cctv.GetAll)
			//cctvGroup.POST("/", cctv.Add)
			//cctvGroup.PUT("/", cctv.Update)
			//cctvGroup.DELETE("/:ID", cctv.Delete)
		}

		organizationGroup := v1.Group("organization")
		{
			organization := handlers.NewOrganizationHandler()
			organizationGroup.GET("/:id", organization.GetOrganizationByID)
			organizationGroup.GET("/", organization.GetAllOrganizations)
			organizationGroup.POST("/", organization.SaveOrganization)
			organizationGroup.PATCH("/", organization.UpdateOrganizationByID)
			organizationGroup.DELETE("/:id", organization.DeleteOrganizationByID)
		}

	}
	return router
}
