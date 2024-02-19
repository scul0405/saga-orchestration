package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/services/order/internal/domain"
	"github.com/scul0405/saga-orchestration/services/order/internal/domain/entity"
	"github.com/scul0405/saga-orchestration/services/order/internal/domain/valueobject"
)

type CreateOrder struct {
	CustomerID uint64
	Products   []PurchasedProduct
}

type PurchasedProduct struct {
	ProductID uint64
	Quantity  uint64
}

type CreateOrderHandler CommandHandler[CreateOrder]

type createOrderHandler struct {
	sf        sonyflake.IDGenerator
	logger    logger.Logger
	orderRepo domain.OrderRepository
}

func NewCreateOrderHandler(sf sonyflake.IDGenerator, logger logger.Logger, orderRepo domain.OrderRepository) CreateOrderHandler {
	return &createOrderHandler{
		sf:        sf,
		logger:    logger,
		orderRepo: orderRepo,
	}
}

func (h *createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) error {
	orderID, err := h.sf.NextID()
	if err != nil {
		return err
	}

	products := make([]valueobject.PurchasedProduct, len(cmd.Products))
	for i, p := range cmd.Products {
		products[i] = valueobject.PurchasedProduct{
			ID:       p.ProductID,
			Quantity: p.Quantity,
		}
	}

	err = h.orderRepo.CreateOrder(ctx, &entity.Order{
		ID:                orderID,
		CustomerID:        cmd.CustomerID,
		PurchasedProducts: &products,
	})

	if err != nil {
		return err
	}

	return nil
}
