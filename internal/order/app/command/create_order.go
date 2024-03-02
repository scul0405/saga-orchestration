package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/order/domain"
	"github.com/scul0405/saga-orchestration/internal/order/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/order/domain/valueobject"
	"github.com/scul0405/saga-orchestration/pkg/logger"
)

type CreateOrder struct {
	OrderID    uint64
	CustomerID uint64
	Products   *[]PurchasedProduct
}

type PurchasedProduct struct {
	ProductID uint64
	Quantity  uint64
}

type CreateOrderHandler CommandHandler[CreateOrder]

type createOrderHandler struct {
	logger    logger.Logger
	orderRepo domain.OrderRepository
}

func NewCreateOrderHandler(logger logger.Logger, orderRepo domain.OrderRepository) CreateOrderHandler {
	return &createOrderHandler{
		logger:    logger,
		orderRepo: orderRepo,
	}
}

func (h *createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) error {
	products := make([]valueobject.PurchasedProduct, len(*cmd.Products))
	for i, p := range *cmd.Products {
		products[i] = valueobject.PurchasedProduct{
			ID:       p.ProductID,
			Quantity: p.Quantity,
		}
	}

	err := h.orderRepo.CreateOrder(ctx, &entity.Order{
		ID:                cmd.OrderID,
		CustomerID:        cmd.CustomerID,
		PurchasedProducts: &products,
	})

	if err != nil {
		return err
	}

	return nil
}
