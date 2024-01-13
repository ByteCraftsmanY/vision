package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"vision/db"
	"vision/dtos"
	"vision/entities"
	"vision/repositories"
	"vision/services"
	"vision/types"
)

type OrganizationHandler struct {
	organizationService services.OrganizationService
	userService         services.UserService
}

func NewOrganizationHandler() *OrganizationHandler {
	session := db.GetSession()

	organizationRepository := repositories.NewOrganizationRepository(session)
	organizationService := services.NewOrganizationService(organizationRepository)

	userRepository := repositories.NewUserRepository(session)
	userService := services.NewUserService(userRepository)

	return &OrganizationHandler{
		organizationService: organizationService,
		userService:         userService,
	}
}

func (o *OrganizationHandler) GetOrganizationByID(ctx *gin.Context) {
	uri := new(dtos.BaseURI)
	if err := ctx.ShouldBindUri(uri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := types.ParseID(uri.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	organization, user := new(entities.Organization), new(entities.User)
	organization, err = o.organizationService.GetOrganizationByID(id)
	if err == nil {
		user, err = o.userService.GetUserByID(organization.AssociatedUserID)
	}
	if errors.Is(err, gocql.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		zap.L().Error("Got error while retrieving user", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"organization": organization, "user": user})
}

func (o *OrganizationHandler) SaveOrganization(ctx *gin.Context) {
	form := new(dtos.OrganizationForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	organization, err := o.organizationService.CreateNewOrganization(
		&entities.Organization{
			Name:             form.Name,
			Contact:          form.Contact,
			Type:             form.Type,
			Address:          form.Address,
			AssociatedUserID: form.AssociatedUserID,
		},
	)

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

func (o *OrganizationHandler) GetAllOrganizations(ctx *gin.Context) {
	form := new(dtos.OrganizationPagination)
	if err := ctx.ShouldBindJSON(form); !strings.EqualFold(ctx.Request.Method, http.MethodGet) && err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//result, err := c.organizationService.GetOrganizations()
	//if err != nil {
	//	zap.L().Error("Got error while retrieving users", zap.Error(err))
	//	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	return
	//}
	//ctx.AbortWithStatusJSON(http.StatusOK, result)
	return
}

func (o *OrganizationHandler) UpdateOrganizationByID(ctx *gin.Context) {
	uri := new(dtos.BaseURI)
	if err := ctx.ShouldBindUri(uri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := types.ParseID(uri.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	form := new(dtos.OrganizationForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := o.organizationService.CreateNewOrganization(&entities.Organization{Base: entities.Base{ID: id}})
	if errors.Is(err, gocql.ErrNotFound) {
		ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		zap.L().Error("Got error while retrieving user", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, user)
}

func (o *OrganizationHandler) DeleteOrganizationByID(ctx *gin.Context) {
	uri := new(dtos.BaseURI)
	if err := ctx.ShouldBindUri(uri); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := types.ParseID(uri.ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = o.organizationService.DeleteOrganizationByID(id)
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
