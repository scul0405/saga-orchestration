package domain

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/payment/domain/entity"
)

type PaymentRepository interface {
	GetPayment(ctx context.Context, paymentID uint64) (*entity.Payment, error)
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	DeletePayment(ctx context.Context, paymentID uint64) error
}
