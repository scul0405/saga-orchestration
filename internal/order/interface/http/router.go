package http

import (
	"github.com/gin-gonic/gin"
	"github.com/scul0405/saga-orchestration/internal/order/app"
	"github.com/scul0405/saga-orchestration/internal/order/app/query"
	"github.com/scul0405/saga-orchestration/internal/order/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/internal/order/interface/http/dto"
	"net/http"
	"strconv"
)

var (
	OkMessage       = "success"
	ErrInvalidID    = "invalid id"
	ErrInvalidJSON  = "invalid json"
	ErrForbidden    = "forbidden"
	ErrInternal     = "internal error"
	ErrInvalidToken = "invalid token"
)

type Router struct {
	app     app.Application
	authSvc grpc.AuthService
}

func NewRouter(app app.Application, authSvc grpc.AuthService) *Router {
	return &Router{
		app:     app,
		authSvc: authSvc,
	}
}

func (r *Router) GetDetailedOrder(c *gin.Context) {
	customerID := r.extractCustomerID(c)
	if customerID == 0 {
		return
	}

	id := c.Param("id")
	orderID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID})
		return
	}

	order, err := r.app.Queries.GetDetailedOrder.Handle(c, query.GetDetailedOrder{OrderID: orderID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
		return
	}

	if order.CustomerID != customerID {
		c.JSON(http.StatusForbidden, gin.H{"error": ErrForbidden})
		return
	}

	resp := &dto.Order{
		OrderID:  orderID,
		Products: make([]dto.Product, len(*(order.PurchasedProducts))),
	}

	for i, p := range *(order.PurchasedProducts) {
		resp.Products[i] = dto.Product{
			ID:          p.ID,
			CategoryID:  p.CategoryID,
			Name:        p.Name,
			BrandName:   p.BrandName,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		}
	}

	c.JSON(http.StatusOK, resp)
}

func (r *Router) extractCustomerID(c *gin.Context) uint64 {
	id, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
		return 0
	}

	return id.(uint64)
}
