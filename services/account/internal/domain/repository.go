//go:generate mockgen -source repository.go -destination ../service/account/mock/service_mock.go -package mock
package domain

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain/valueobject"
)

// CustomerRepository is the customer repository interface
type CustomerRepository interface {
	GetCustomerPersonalInfo(ctx context.Context, customerID uint64) (*valueobject.CustomerPersonalInfo, error)
	GetCustomerDeliveryInfo(ctx context.Context, customerID uint64) (*valueobject.CustomerDeliveryInfo, error)
	UpdateCustomerPersonalInfo(ctx context.Context, customerID uint64, personalInfo *valueobject.CustomerPersonalInfo) error
	UpdateCustomerDeliveryInfo(ctx context.Context, customerID uint64, deliveryInfo *valueobject.CustomerDeliveryInfo) error
}
