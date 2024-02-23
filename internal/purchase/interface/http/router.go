package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/scul0405/saga-orchestration/internal/purchase/app"
	"github.com/scul0405/saga-orchestration/internal/purchase/app/command"
	"github.com/scul0405/saga-orchestration/internal/purchase/app/query"
	"github.com/scul0405/saga-orchestration/internal/purchase/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/internal/purchase/interface/http/dto"
	"net/http"
)

var (
	OkMessage          = "success"
	ErrInvalidID       = "invalid id"
	ErrInvalidJSON     = "invalid json"
	ErrForbidden       = "forbidden"
	ErrInternal        = "internal error"
	ErrInvalidToken    = "invalid token"
	ErrProductNotFound = "product not found"
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

func (r *Router) CreatePurchase(c *gin.Context) {
	customerID := r.extractCustomerID(c)
	if customerID == 0 {
		return
	}

	var purchase dto.Purchase
	if err := c.ShouldBindJSON(&purchase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	orderItemsQuery := make([]query.OrderItem, len(*purchase.OrderItems))
	for i, item := range *purchase.OrderItems {
		orderItemsQuery[i] = query.OrderItem{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		}
	}
	checkQuery := query.CheckProducts{
		OrderItems: &orderItemsQuery,
	}
	productStatuses, err := r.app.Queries.CheckProducts.Handle(c, checkQuery)
	if err != nil {
		if errors.Is(err, query.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrProductNotFound})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
			return
		}
	}

	purchaseCmd := command.CreatePurchase{
		Order: &command.Order{
			CustomerID: customerID,
		},
		Payment: &command.Payment{
			CurrencyCode: purchase.Payment.CurrencyCode,
		},
	}

	var amount uint64
	for i, item := range *productStatuses {
		amount += (*purchase.OrderItems)[i].Quantity * item.Price
	}
	purchaseCmd.Payment.Amount = amount

	orderItemsCmd := make([]command.OrderItem, len(*purchase.OrderItems))
	for i, item := range *purchase.OrderItems {
		orderItemsCmd[i] = command.OrderItem{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		}
	}
	purchaseCmd.Order.OrderItems = &orderItemsCmd

	err = r.app.Commands.CreatePurchase.Handle(c, purchaseCmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
		return
	}
}

func (r *Router) extractCustomerID(c *gin.Context) uint64 {
	id, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidToken})
		return 0
	}

	return id.(uint64)
}
