package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"vision/forms"
	"vision/services"
)

type OrgController struct{}

var orgService = new(services.OrgService)

func (c *OrgController) Add(ctx *gin.Context) {
	form := new(forms.OrganizationNew)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	organization, err := orgService.Add(form)
	if err != nil && strings.Contains(err.Error(), "exists") {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		zap.L().Error("Got error while registering organization", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusCreated, organization)
}

func (c *OrgController) Get(ctx *gin.Context) {
	id := ctx.Param("ID")
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	organization, err := orgService.GetByID(uuid)
	if errors.Is(err, gocql.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		zap.L().Error("Got error while retrieving user", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, organization)
}

func (c *OrgController) GetAll(ctx *gin.Context) {
	form := new(forms.OrganizationPagination)
	if err := ctx.ShouldBindJSON(form); !strings.EqualFold(ctx.Request.Method, http.MethodGet) && err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := orgService.GetAll(form)
	if err != nil {
		zap.L().Error("Got error while retrieving users", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, result)
	return
}

func (c *OrgController) Update(ctx *gin.Context) {
	form := new(forms.OrganizationEdit)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := orgService.Update(form)
	if errors.Is(err, gocql.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		zap.L().Error("Got error while retrieving user", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, user)
}

func (c *OrgController) Delete(ctx *gin.Context) {
	id := ctx.Param("ID")
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = orgService.Remove(uuid)
	if errors.Is(err, gocql.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		zap.L().Error("Got error while retrieving user", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatus(http.StatusNoContent)
}
