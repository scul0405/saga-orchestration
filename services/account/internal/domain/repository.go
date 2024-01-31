package domain

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain/valueobject"
)

// CustomerRepository is the customer repository interface
type CustomerRepository interface {
	GetCustomerPersonalInfo(ctx context.Context, customerID uint64) (*CustomerPersonalInfo, error)
	GetCustomerDeliveryInfo(ctx context.Context, customerID uint64) (*CustomerDeliveryInfo, error)
	UpdateCustomerPersonalInfo(ctx context.Context, customerID uint64, personalInfo *valueobject.CustomerPersonalInfo) error
	UpdateCustomerShippingInfo(ctx context.Context, customerID uint64, shippingInfo *valueobject.CustomerDeliveryInfo) error
}

// CustomerPersonalInfo os customer personal info type
type CustomerPersonalInfo struct {
	FirstName string
	LastName  string
	Email     string
}

// CustomerDeliveryInfo os customer shipping info type
type CustomerDeliveryInfo struct {
	Address     string
	PhoneNumber string
}
