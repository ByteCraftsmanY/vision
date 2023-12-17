package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"net/http"
	"vision/forms"
)

type AuthController struct{}

func (c AuthController) Login(ctx *gin.Context) {
	form := new(forms.Login)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := userModel.LoginCheck(form)
	if errors.Is(err, gocql.ErrNotFound) {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	} else if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"token": token})
}
