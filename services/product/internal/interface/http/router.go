package http

import (
	"github.com/gin-gonic/gin"
	"github.com/scul0405/saga-orchestration/services/product/internal/app"
	"github.com/scul0405/saga-orchestration/services/product/internal/app/command"
	"github.com/scul0405/saga-orchestration/services/product/internal/app/query"
	"github.com/scul0405/saga-orchestration/services/product/internal/interface/http/dto"
	"github.com/scul0405/saga-orchestration/services/product/internal/service"
	"net/http"
	"strconv"
)

var (
	OkMessage       = "success"
	ErrInvalidID    = "invalid id"
	ErrInvalidJSON  = "invalid json"
	ErrInternal     = "internal error"
	ErrInvalidToken = "invalid token"
)

type Router struct {
	app     app.Application
	authSvc service.AuthService
}

func NewRouter(app app.Application, authSvc service.AuthService) *Router {
	return &Router{
		app:     app,
		authSvc: authSvc,
	}
}

func (r *Router) CreateProduct(c *gin.Context) {
	var product dto.CreateProduct
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	err := r.app.Commands.CreateProduct.Handle(c, command.CreateProduct{
		CategoryID:  product.CategoryID,
		Name:        product.Name,
		BrandName:   product.BrandName,
		Description: product.Description,
		Price:       product.Price,
		Inventory:   product.Inventory,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": OkMessage})
}

func (r *Router) UpdateProductDetail(c *gin.Context) {
	var product dto.UpdateProductDetail
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	err := r.app.Commands.UpdateProductDetail.Handle(c, command.UpdateProductDetail{
		ProductID:   product.ID,
		Name:        product.Name,
		BrandName:   product.BrandName,
		Description: product.Description,
		Price:       product.Price,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": OkMessage})
}

func (r *Router) GetProduct(c *gin.Context) {
	idParam := c.Param("id")

	productID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID})
		return
	}

	productsID := []uint64{uint64(productID)}

	products, err := r.app.Queries.GetProducts.Handle(c, query.GetProducts{ProductIDs: &productsID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
		return
	}

	product := (*products)[0]

	c.JSON(http.StatusOK, &dto.Product{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Name:        product.Detail.Name,
		BrandName:   product.Detail.BrandName,
		Description: product.Detail.Description,
		Price:       product.Detail.Price,
		Inventory:   product.Inventory,
	})
}

func (r *Router) extractCustomerID(c *gin.Context) uint64 {
	id, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
		return 0
	}

	return id.(uint64)
}