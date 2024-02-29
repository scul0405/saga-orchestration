package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/aggregate"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/valueobject"
	"github.com/scul0405/saga-orchestration/internal/purchase/eventhandler"
	"github.com/scul0405/saga-orchestration/internal/purchase/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
)

type CreatePurchase struct {
	Order   *Order
	Payment *Payment
}

type Order struct {
	CustomerID uint64
	OrderItems *[]OrderItem
}

type OrderItem struct {
	ID       uint64
	Price    uint64
	Quantity uint64
}

type Payment struct {
	CurrencyCode string
	Amount       uint64
}

type CreatePurchaseHandler CommandHandler[CreatePurchase]

type createPurchaseHandler struct {
	sf         sonyflake.IDGenerator
	logger     logger.Logger
	productSvc grpc.ProductService
	evPub      eventhandler.PurchaseEventHandler
}

func NewCreatePurchaseHandler(
	sf sonyflake.IDGenerator, logger logger.Logger,
	productSvc grpc.ProductService,
	evPub eventhandler.PurchaseEventHandler) CreatePurchaseHandler {
	return &createPurchaseHandler{
		sf:         sf,
		logger:     logger,
		productSvc: productSvc,
		evPub:      evPub,
	}
}

func (h *createPurchaseHandler) Handle(ctx context.Context, cmd CreatePurchase) error {
	orderItems := make([]entity.OrderItem, len(*cmd.Order.OrderItems))
	for i, item := range *cmd.Order.OrderItems {
		orderItems[i] = entity.OrderItem{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	purchaseID, err := h.sf.NextID()
	if err != nil {
		return err
	}

	aggOrderItems := make([]entity.OrderItem, len(*cmd.Order.OrderItems))
	for i, item := range *cmd.Order.OrderItems {
		aggOrderItems[i] = entity.OrderItem{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	purchase := &aggregate.Purchase{
		ID: purchaseID,
		Order: &entity.Order{
			CustomerID: cmd.Order.CustomerID,
			OrderItems: &aggOrderItems,
		},
		Payment: &valueobject.Payment{
			CurrencyCode: cmd.Payment.CurrencyCode,
			Amount:       cmd.Payment.Amount,
		},
	}

	return h.evPub.CreatePurchase(ctx, purchase)
}