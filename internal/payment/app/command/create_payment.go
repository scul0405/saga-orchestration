package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/payment/domain"
	"github.com/scul0405/saga-orchestration/internal/payment/domain/entity"
	"github.com/scul0405/saga-orchestration/pkg/logger"
)

type CreatePayment struct {
	ID           uint64
	CustomerID   uint64
	Amount       uint64
	CurrencyCode string
}

type CreatePaymentHandler CommandHandler[CreatePayment]

type createPaymentHandler struct {
	logger      logger.Logger
	paymentRepo domain.PaymentRepository
}

func NewCreatePaymentHandler(logger logger.Logger, paymentRepo domain.PaymentRepository) CreatePaymentHandler {
	return &createPaymentHandler{
		logger:      logger,
		paymentRepo: paymentRepo,
	}
}

func (h *createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) error {
	err := h.paymentRepo.CreatePayment(ctx, &entity.Payment{
		ID:           cmd.ID,
		CustomerID:   cmd.CustomerID,
		Amount:       cmd.Amount,
		CurrencyCode: cmd.CurrencyCode,
	})

	if err != nil {
		return err
	}

	return nil
}
