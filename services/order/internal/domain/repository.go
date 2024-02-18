package domain

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/order/internal/domain/entity"
)

type OrderRepository interface {
	GetOrder(ctx context.Context, id uint64) (*entity.Order, error)
	CreateOrder(ctx context.Context, order *entity.Order) error
	DeleteOrder(ctx context.Context, id uint64) error
}
