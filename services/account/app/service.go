package app

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/account/domain/entity"
	valueobject2 "github.com/scul0405/saga-orchestration/services/account/domain/valueobject"
)

// CustomerService is the service for customer domain
type CustomerService interface {
	GetPersonalInfo(ctx context.Context, customerID uint64) (*valueobject2.CustomerPersonalInfo, error)
	GetDeliveryInfo(ctx context.Context, customerID uint64) (*valueobject2.CustomerDeliveryInfo, error)
	UpdatePersonalInfo(ctx context.Context, customerID uint64, info *valueobject2.CustomerPersonalInfo) error
	UpdateDeliveryInfo(ctx context.Context, customerID uint64, info *valueobject2.CustomerDeliveryInfo) error
}

// AuthService is the service for authentication
type AuthService interface {
	Auth(ctx context.Context, authPayload *valueobject2.AuthPayload) (*valueobject2.AuthResponse, error)
	Register(ctx context.Context, customer *entity.Customer) (string, string, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, error)
}
