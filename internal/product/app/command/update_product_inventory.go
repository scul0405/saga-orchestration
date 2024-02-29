package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/product/domain"
	"github.com/scul0405/saga-orchestration/internal/product/domain/valueobject"
	"github.com/scul0405/saga-orchestration/pkg/logger"
)

type UpdateProductInventory struct {
	IdempotencyKey    uint64
	PurchasedProducts *[]PurchasedProduct
}

type PurchasedProduct struct {
	ID       uint64
	Quantity uint64
}

type UpdateProductInventoryHandler CommandHandler[UpdateProductInventory]

type updateProductInventoryHandler struct {
	logger      logger.Logger
	productRepo domain.ProductRepository
}

func NewUpdateProductInventoryHandler(logger logger.Logger, productRepo domain.ProductRepository) UpdateProductInventoryHandler {
	return &updateProductInventoryHandler{
		logger:      logger,
		productRepo: productRepo,
	}
}

func (h *updateProductInventoryHandler) Handle(ctx context.Context, cmd UpdateProductInventory) error {
	purchaseProducts := make([]valueobject.PurchasedProduct, len(*cmd.PurchasedProducts))
	for i, p := range *cmd.PurchasedProducts {
		purchaseProducts[i] = valueobject.PurchasedProduct{
			ID:       p.ID,
			Quantity: p.Quantity,
		}
	}
	err := h.productRepo.UpdateProductInventory(ctx, cmd.IdempotencyKey, &purchaseProducts)

	if err != nil {
		return err
	}

	return nil
}
