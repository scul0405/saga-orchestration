package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/payment/internal/domain"
	"github.com/scul0405/saga-orchestration/services/payment/internal/domain/entity"
	"github.com/scul0405/saga-orchestration/services/payment/internal/infrastructure/logger"
	"github.com/scul0405/saga-orchestration/services/payment/pkg"
)

type CreatePayment struct {
	CustomerID   uint64
	Amount       uint64
	CurrencyCode string
}

type CreatePaymentHandler CommandHandler[CreatePayment]

type createPaymentHandler struct {
	sf          pkg.IDGenerator
	logger      logger.Logger
	paymentRepo domain.PaymentRepository
}

func NewCreatePaymentHandler(sf pkg.IDGenerator, logger logger.Logger, paymentRepo domain.PaymentRepository) CreatePaymentHandler {
	return &createPaymentHandler{
		sf:          sf,
		logger:      logger,
		paymentRepo: paymentRepo,
	}
}

func (h *createPaymentHandler) Handle(ctx context.Context, cmd CreatePayment) error {
	orderID, err := h.sf.NextID()
	if err != nil {
		return err
	}

	err = h.paymentRepo.CreatePayment(ctx, &entity.Payment{
		ID:           orderID,
		CustomerID:   cmd.CustomerID,
		Amount:       cmd.Amount,
		CurrencyCode: cmd.CurrencyCode,
	})

	if err != nil {
		return err
	}

	return nil
}
