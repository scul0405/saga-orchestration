package http

import (
	"github.com/gin-gonic/gin"
	"github.com/scul0405/saga-orchestration/services/payment/internal/app"
	"github.com/scul0405/saga-orchestration/services/payment/internal/app/query"
	"github.com/scul0405/saga-orchestration/services/payment/internal/infrastructure/grpc/auth"
	"github.com/scul0405/saga-orchestration/services/payment/internal/interface/http/dto"
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
	authSvc auth.AuthService
}

func NewRouter(app app.Application, authSvc auth.AuthService) *Router {
	return &Router{
		app:     app,
		authSvc: authSvc,
	}
}

func (r *Router) GetPayment(c *gin.Context) {
	customerID := r.extractCustomerID(c)
	if customerID == 0 {
		return
	}

	id := c.Param("id")
	paymentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidID})
		return
	}

	payment, err := r.app.Queries.GetPayment.Handle(c, query.GetPayment{PaymentID: paymentID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
		return
	}

	if payment.CustomerID != customerID {
		c.JSON(http.StatusForbidden, gin.H{"error": ErrForbidden})
		return
	}

	c.JSON(http.StatusOK, &dto.Payment{
		ID:           payment.ID,
		CustomerID:   payment.CustomerID,
		Amount:       payment.Amount,
		CurrencyCode: payment.CurrencyCode,
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
