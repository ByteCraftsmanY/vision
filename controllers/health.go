package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vision/models"
)

type HealthController struct{}

var healthModel = new(models.Health)

func (h HealthController) Status(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusOK, healthModel.Get())
}
