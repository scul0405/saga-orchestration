package app

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/account/internal/domain/valueobject"
)

// CustomerService is the service for customer domain
type CustomerService interface {
	GetPersonalInfo(ctx context.Context, customerID uint64) (*valueobject.CustomerPersonalInfo, error)
	GetDeliveryInfo(ctx context.Context, customerID uint64) (*valueobject.CustomerDeliveryInfo, error)
	UpdatePersonalInfo(ctx context.Context, customerID uint64, info *valueobject.CustomerPersonalInfo) error
	UpdateDeliveryInfo(ctx context.Context, customerID uint64, info *valueobject.CustomerDeliveryInfo) error
}
