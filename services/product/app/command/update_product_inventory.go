package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/services/product/domain"
	"github.com/scul0405/saga-orchestration/services/product/domain/valueobject"
)

type UpdateProductInventory struct {
	IdempotencyKey    uint64
	PurchasedProducts *[]valueobject.PurchasedProduct
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
	err := h.productRepo.UpdateProductInventory(ctx, cmd.IdempotencyKey, cmd.PurchasedProducts)

	if err != nil {
		return err
	}

	return nil
}
