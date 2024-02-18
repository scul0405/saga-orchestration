package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/services/order/internal/domain"
	"github.com/scul0405/saga-orchestration/services/order/internal/infrastructure/logger"
)

type DeleteOrder struct {
	OrderID uint64
}

type DeleteOrderHandler CommandHandler[DeleteOrder]

type deleteOrderHandler struct {
	logger    logger.Logger
	orderRepo domain.OrderRepository
}

func NewDeleteOrderHandler(logger logger.Logger, orderRepo domain.OrderRepository) DeleteOrderHandler {
	return &deleteOrderHandler{
		logger:    logger,
		orderRepo: orderRepo,
	}
}

func (h *deleteOrderHandler) Handle(ctx context.Context, cmd DeleteOrder) error {
	err := h.orderRepo.DeleteOrder(ctx, cmd.OrderID)
	if err != nil {
		return err
	}

	return nil
}
