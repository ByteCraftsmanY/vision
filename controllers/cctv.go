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

type CCTVController struct{}

var cctvService = new(services.CCTVService)

func (c *CCTVController) Add(ctx *gin.Context) {
	form := new(forms.CCTV)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cctv, err := cctvService.Add(form)
	if err != nil && strings.Contains(err.Error(), "exists") {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		zap.L().Error("Got error while registering cctv", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusCreated, cctv)
}

func (c *CCTVController) Get(ctx *gin.Context) {
	id := ctx.Param("ID")
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := cctvService.GetByID(uuid)
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

func (c *CCTVController) GetAll(ctx *gin.Context) {
	form := new(forms.CCTVPagination)
	if err := ctx.ShouldBindJSON(form); !strings.EqualFold(ctx.Request.Method, http.MethodGet) && err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := cctvService.GetAll(form)
	if err != nil {
		zap.L().Error("Got error while retrieving users", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, result)
	return
}

func (c *CCTVController) Update(ctx *gin.Context) {
	form := new(forms.CCTVEdit)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := cctvService.Update(form)
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

func (c *CCTVController) Delete(ctx *gin.Context) {
	id := ctx.Param("ID")
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = cctvService.Remove(uuid)
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
