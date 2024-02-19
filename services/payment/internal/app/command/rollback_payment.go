package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/services/payment/internal/domain"
)

type RollbackPayment struct {
	PaymentID uint64
}

type RollbackPaymentHandler CommandHandler[RollbackPayment]

type rollbackPaymentHandler struct {
	logger      logger.Logger
	paymentRepo domain.PaymentRepository
}

func NewRollbackPaymentHandler(logger logger.Logger, paymentRepo domain.PaymentRepository) RollbackPaymentHandler {
	return &rollbackPaymentHandler{
		logger:      logger,
		paymentRepo: paymentRepo,
	}
}

func (h *rollbackPaymentHandler) Handle(ctx context.Context, cmd RollbackPayment) error {
	err := h.paymentRepo.DeletePayment(ctx, cmd.PaymentID)
	if err != nil {
		return err
	}

	return nil
}
