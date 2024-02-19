package query

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/payment/internal/domain"
	"github.com/scul0405/saga-orchestration/services/payment/internal/domain/entity"
	"github.com/scul0405/saga-orchestration/services/payment/internal/infrastructure/logger"
)

type GetPayment struct {
	PaymentID uint64
}

type GetPaymentHandler QueryHandler[GetPayment, *entity.Payment]

type getPaymentHandler struct {
	logger      logger.Logger
	paymentRepo domain.PaymentRepository
}

func NewGetPaymentHandler(logger logger.Logger, paymentRepo domain.PaymentRepository) GetPaymentHandler {
	return &getPaymentHandler{
		logger:      logger,
		paymentRepo: paymentRepo,
	}
}

func (h *getPaymentHandler) Handle(ctx context.Context, query GetPayment) (*entity.Payment, error) {
	payment, err := h.paymentRepo.GetPayment(ctx, query.PaymentID)
	if err != nil {
		return nil, err
	}

	return &entity.Payment{
		ID:           payment.ID,
		CustomerID:   payment.CustomerID,
		Amount:       payment.Amount,
		CurrencyCode: payment.CurrencyCode,
	}, nil
}
