package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"vision/constants"
	"vision/db"
	"vision/dtos"
	"vision/repositories"
	"vision/services"
	"vision/utils/token"
)

type AuthController struct {
	AuthService services.AuthService
	UserService services.UserService
}

func NewAuthController() *AuthController {
	session := db.GetSession()

	authRepository := repositories.NewAuthRepository(session)
	userRepository := repositories.NewUserRepository(session)

	authService := services.NewAuthService(authRepository)
	userService := services.NewUserService(userRepository)
	return &AuthController{
		AuthService: authService,
		UserService: userService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	form := new(dtos.Login)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.UserService.GetUserByEmail(form.Email)
	if errors.Is(err, gocql.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		zap.L().Error("Got error while retrieving user", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authToken, err := token.Generate(user.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(
		http.StatusOK,
		gin.H{
			"status": "success",
			"token":  authToken,
			"user":   user,
		},
	)
}

func (c *AuthController) GenerateOTP(ctx *gin.Context) {
	form := new(dtos.OTPCreateRequest)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	auth, err := c.AuthService.CreateNewAuthWithTTL(form, constants.AuthOtpTTL)
	if err != nil && strings.Contains(err.Error(), "exists") {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		zap.L().Error("Got error while generating auth", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusCreated, auth)
}

func (c *AuthController) VerifyOTP(ctx *gin.Context) {
	form := new(dtos.OTPVerifyRequest)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auth, err := c.AuthService.GetAuthByPhone(form.Phone)
	if errors.Is(err, gocql.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		zap.L().Error("Got error while retrieving auth", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if auth.Code != form.Code {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"error": "Wrong Code"})
		return
	}

	user, err := c.UserService.GetUserByPhone(auth.Phone)
	if errors.Is(err, gocql.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"status": "success", "msg": "create new account"})
		return
	}
	if err != nil {
		zap.L().Error("Got error while retrieving user", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	authToken, err := token.Generate(user.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(
		http.StatusOK,
		gin.H{
			"status": "success",
			"token":  authToken,
			"auth":   auth,
			"user":   user,
		},
	)
}
