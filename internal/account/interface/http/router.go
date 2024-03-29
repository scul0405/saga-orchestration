package http

import (
	"github.com/gin-gonic/gin"
	"github.com/scul0405/saga-orchestration/internal/account/app"
	"github.com/scul0405/saga-orchestration/internal/account/interface/http/dto"
	"github.com/scul0405/saga-orchestration/internal/account/repository/postgres_repo"
	"github.com/scul0405/saga-orchestration/internal/account/service/auth"
	"net/http"
)

var (
	OkMessage      = "success"
	ErrInvalidJSON = "invalid json"
	ErrInternal    = "internal error"
)

type Router struct {
	authSvc     app.AuthService
	customerSvc app.CustomerService
}

func NewRouter(authSvc app.AuthService, customerSvc app.CustomerService) *Router {
	return &Router{
		authSvc:     authSvc,
		customerSvc: customerSvc,
	}
}

func (r *Router) Register(c *gin.Context) {
	var customer dto.RegisterCustomer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	domainObject := customer.ToDomainObject()
	accessToken, refreshToken, err := r.authSvc.Register(c, &domainObject)

	switch err {
	case nil:
		c.JSON(http.StatusCreated, &dto.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken})
	case postgres_repo.ErrDuplicateEntry:
		c.JSON(http.StatusConflict, gin.H{"error": postgres_repo.ErrDuplicateEntry})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
	}
}

func (r *Router) Login(c *gin.Context) {
	var customer dto.LoginCustomer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	accessToken, refreshToken, err := r.authSvc.Login(c, customer.Email, customer.Password)

	switch err {
	case nil:
		c.JSON(http.StatusOK, &dto.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken})
	case auth.ErrAuthenticationFailed:
		c.JSON(http.StatusUnauthorized, gin.H{"error": auth.ErrAuthenticationFailed})
	case auth.ErrCustomerInactive:
		c.JSON(http.StatusForbidden, gin.H{"error": auth.ErrCustomerInactive})
	case auth.ErrCustomerNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": auth.ErrCustomerNotFound})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
	}
}

func (r *Router) RefreshToken(c *gin.Context) {
	var token dto.RefreshToken
	if err := c.ShouldBindJSON(&token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	accessToken, refreshToken, err := r.authSvc.RefreshToken(c, token.RefreshToken)

	switch err {
	case nil:
		c.JSON(http.StatusOK, &dto.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken})
	case auth.ErrInvalidToken:
		c.JSON(http.StatusUnauthorized, gin.H{"error": auth.ErrInvalidToken})
	case auth.ErrTokenExpired:
		c.JSON(http.StatusUnauthorized, gin.H{"error": auth.ErrTokenExpired})
	case auth.ErrCustomerInactive:
		c.JSON(http.StatusForbidden, gin.H{"error": auth.ErrCustomerInactive})
	case auth.ErrCustomerNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": auth.ErrCustomerNotFound})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
	}
}

func (r *Router) GetCustomerPersonalInfo(c *gin.Context) {
	customerID := r.extractCustomerID(c)
	info, err := r.customerSvc.GetPersonalInfo(c, customerID)

	switch err {
	case nil:
		c.JSON(http.StatusOK, &dto.CustomerPersonalInfo{
			FirstName: info.FirstName,
			LastName:  info.LastName,
			Email:     info.Email,
		})
	case auth.ErrCustomerNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": auth.ErrCustomerNotFound})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
	}
}

func (r *Router) UpdateCustomerPersonalInfo(c *gin.Context) {
	var info dto.CustomerPersonalInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	customerID := r.extractCustomerID(c)
	err := r.customerSvc.UpdatePersonalInfo(c, customerID, info.ToDomainObject())

	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"message": OkMessage})
	case postgres_repo.ErrCustomerNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": postgres_repo.ErrCustomerNotFound})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
	}
}

func (r *Router) GetCustomerDeliveryInfo(c *gin.Context) {
	customerID := r.extractCustomerID(c)
	info, err := r.customerSvc.GetDeliveryInfo(c, customerID)

	switch err {
	case nil:
		c.JSON(http.StatusOK, &dto.CustomerDeliveryInfo{
			Address:     info.Address,
			PhoneNumber: info.PhoneNumber,
		})
	case auth.ErrCustomerNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": auth.ErrCustomerNotFound})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
	}
}

func (r *Router) UpdateCustomerDeliveryInfo(c *gin.Context) {
	var info dto.CustomerDeliveryInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON})
		return
	}

	customerID := r.extractCustomerID(c)
	err := r.customerSvc.UpdateDeliveryInfo(c, customerID, info.ToDomainObject())

	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{"message": OkMessage})
	case postgres_repo.ErrCustomerNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": postgres_repo.ErrCustomerNotFound})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrInternal})
	}
}

func (r *Router) extractCustomerID(c *gin.Context) uint64 {
	id, exists := c.Get("customer_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": auth.ErrInvalidToken})
		return 0
	}

	return id.(uint64)
}
