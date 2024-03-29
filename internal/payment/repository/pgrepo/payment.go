package pgrepo

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/payment/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/payment/infrastructure/db/postgres/model"
	"gorm.io/gorm"
)

type paymentRepositoryImpl struct {
	db *gorm.DB
}

type PaymentRepository interface {
	GetPayment(ctx context.Context, paymentID uint64) (*entity.Payment, error)
	CreatePayment(ctx context.Context, payment *entity.Payment) error
	DeletePayment(ctx context.Context, paymentID uint64) error
}

func NewOrderRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepositoryImpl{db: db}
}

func (r *paymentRepositoryImpl) GetPayment(ctx context.Context, paymentID uint64) (*entity.Payment, error) {
	var payment model.Payment
	if err := r.db.Where("id = ?", paymentID).First(&payment).WithContext(ctx).Error; err != nil {
		return nil, err
	}

	return &entity.Payment{
		ID:           payment.ID,
		CustomerID:   payment.CustomerID,
		Amount:       payment.Amount,
		CurrencyCode: payment.CurrencyCode,
	}, nil
}

func (r *paymentRepositoryImpl) CreatePayment(ctx context.Context, payment *entity.Payment) error {
	paymentModel := model.Payment{
		ID:           payment.ID,
		CustomerID:   payment.CustomerID,
		Amount:       payment.Amount,
		CurrencyCode: payment.CurrencyCode,
	}

	if err := r.db.Create(&paymentModel).WithContext(ctx).Error; err != nil {
		return err
	}

	return nil
}

func (r *paymentRepositoryImpl) DeletePayment(ctx context.Context, paymentID uint64) error {
	if err := r.db.Exec("DELETE FROM payments WHERE id = ?", paymentID).WithContext(ctx).Error; err != nil {
		return err
	}
	return nil
}
