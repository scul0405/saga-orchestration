package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/product/domain"
	"github.com/scul0405/saga-orchestration/internal/product/domain/valueobject"
	"github.com/scul0405/saga-orchestration/pkg/logger"
)

type RollbackProductInventory struct {
	IdempotencyKey    uint64
	PurchasedProducts *[]PurchasedProduct
}

type PurchaseProduct struct {
	ID       uint64
	Quantity uint64
}

type RollbackProductInventoryHandler CommandHandler[RollbackProductInventory]

type rollbackProductInventoryHandler struct {
	logger      logger.Logger
	productRepo domain.ProductRepository
}

func NewRollbackProductInventoryHandler(logger logger.Logger, productRepo domain.ProductRepository) RollbackProductInventoryHandler {
	return &rollbackProductInventoryHandler{
		logger:      logger,
		productRepo: productRepo,
	}
}

func (h *rollbackProductInventoryHandler) Handle(ctx context.Context, cmd RollbackProductInventory) error {
	purchasedProducts := make([]valueobject.PurchasedProduct, len(*cmd.PurchasedProducts))
	for i, item := range *cmd.PurchasedProducts {
		purchasedProducts[i] = valueobject.PurchasedProduct{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}
	err := h.productRepo.RollbackProductInventory(ctx, cmd.IdempotencyKey, &purchasedProducts)

	if err != nil {
		return err
	}

	return nil
}
