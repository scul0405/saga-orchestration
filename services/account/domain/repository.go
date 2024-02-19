//go:generate mockgen -source repository.go -destination ../service/mock/repository_mock.go -package mock
package domain

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/account/domain/entity"
	valueobject2 "github.com/scul0405/saga-orchestration/services/account/domain/valueobject"
)

// CustomerRepository is the customer repository interface
type CustomerRepository interface {
	GetCustomerPersonalInfo(ctx context.Context, customerID uint64) (*valueobject2.CustomerPersonalInfo, error)
	GetCustomerDeliveryInfo(ctx context.Context, customerID uint64) (*valueobject2.CustomerDeliveryInfo, error)
	UpdateCustomerPersonalInfo(ctx context.Context, customerID uint64, personalInfo *valueobject2.CustomerPersonalInfo) error
	UpdateCustomerDeliveryInfo(ctx context.Context, customerID uint64, deliveryInfo *valueobject2.CustomerDeliveryInfo) error
}

// JWTAuthRepository is the jwt auth repository interface
type JWTAuthRepository interface {
	// CheckCustomer checks if the customer exists and is active
	CheckCustomer(ctx context.Context, customerID uint64) (bool, bool, error)
	CreateCustomer(ctx context.Context, customer *entity.Customer) error
	GetCustomerCredentials(ctx context.Context, email string) (bool, *valueobject2.CustomerCredentials, error)
}
