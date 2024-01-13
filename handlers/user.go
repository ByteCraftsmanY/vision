package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"sync"
	"time"
	"vision/db"
	"vision/dtos"
	"vision/entities"
	"vision/repositories"
	"vision/services"
	"vision/utils/token"
)

type UserController struct {
	UserService         services.UserService
	OrganizationService services.OrganizationService
}

func NewUserController() *UserController {
	userRepository := repositories.NewUserRepository(db.GetSession())
	organizationRepository := repositories.NewOrganizationRepository(db.GetSession())

	userService := services.NewUserService(userRepository)
	organizationService := services.NewOrganizationService(organizationRepository)
	return &UserController{
		UserService:         userService,
		OrganizationService: organizationService,
	}
}

func (c UserController) Retrieve(ctx *gin.Context) {
	id, err := token.ExtractUserID(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := c.UserService.GetUserByID(id)
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

func (c UserController) Store(ctx *gin.Context) {
	form := new(dtos.UserForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := isUserExists(c.UserService, form)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := c.UserService.CreateNewUser(
		&entities.User{
			Password: form.Password,
			Name:     form.Name,
			Email:    form.Email,
			Phone:    form.Phone,
		},
	)
	if err != nil && strings.Contains(err.Error(), "exists") {
		ctx.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		zap.L().Error("Got error while storing user", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusCreated, user)
}

func isUserExists(s services.UserService, form *dtos.UserForm) error {
	type userErr struct {
		Type string
		User *entities.User
		Err  error
	}

	wg := sync.WaitGroup{}
	userChan := make(chan userErr)
	wg.Add(2)
	go func(userChan chan userErr, form *dtos.UserForm, wg *sync.WaitGroup) {
		defer wg.Done()
		user, err := s.GetUserByEmail(form.Email)
		userChan <- userErr{
			Type: "email",
			User: user,
			Err:  err,
		}
	}(userChan, form, &wg)
	go func(userChan chan userErr, form *dtos.UserForm, wg *sync.WaitGroup) {
		defer wg.Done()
		user, err := s.GetUserByPhone(form.Phone)
		time.Sleep(time.Second * 5)
		userChan <- userErr{
			Type: "phone",
			User: user,
			Err:  err,
		}
	}(userChan, form, &wg)
	go func(userChan chan userErr, wg *sync.WaitGroup) {
		wg.Wait()
		close(userChan)
	}(userChan, &wg)

	for ue := range userChan {
		fmt.Printf("time-%v\ntype %v\nval %+v\nerr %v\n", time.Now(), ue.Type, ue.User, ue.Err)
		if ue.Err != nil {
			return fmt.Errorf("failed to fetch user data by %s", ue.Type)
		}
		if len(ue.User.Name) > 0 {
			return fmt.Errorf("user is already exists with this %s", ue.Type)
		}
	}
	return nil
}
