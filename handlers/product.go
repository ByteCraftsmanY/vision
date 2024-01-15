package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"vision/daos"
	"vision/db"
	"vision/dtos"
	"vision/repositories"
	"vision/services"
	"vision/types"
)

type ProductHandler struct {
	ProductService services.ProductService
}

func NewProductHandler() *ProductHandler {
	productRepo := repositories.NewProductRepository(db.GetSession())
	productService := services.NewProductService(productRepo)
	return &ProductHandler{
		ProductService: productService,
	}
}

func (c *ProductHandler) Add(ctx *gin.Context) {
	form := new(dtos.ProductForm)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cctv, err := c.ProductService.CreateNewProduct(
		&daos.Product{
			Username:       form.UserName,
			Password:       form.Password,
			URL:            form.URL,
			OrganizationID: form.OrganizationID,
			Base:           daos.Base{Extra: form.Extra},
		},
	)
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

func (c *ProductHandler) Get(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := types.ParseID(idStr)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := c.ProductService.GetProductByID(id)
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

/*
func (c *ProductHandler) GetAll(ctx *gin.Context) {
	form := new(dtos.CCTVPagination)
	if err := ctx.ShouldBindJSON(form); !strings.EqualFold(ctx.Request.Method, http.MethodGet) && err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := c.ProductService.GetProducts(form)
	if err != nil {
		zap.L().Error("Got error while retrieving users", zap.Error(err))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.AbortWithStatusJSON(http.StatusOK, result)
	return
}

func (c *ProductHandler) Update(ctx *gin.Context) {
	form := new(dtos.CCTVEdit)
	if err := ctx.ShouldBindJSON(form); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := c.ProductService.Update(form)
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

func (c *ProductHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("ID")
	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = c.ProductService.Remove(uuid)
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
*/
