package command

import (
	"context"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/services/product/domain"
	"github.com/scul0405/saga-orchestration/services/product/domain/valueobject"
)

type UpdateProductDetail struct {
	ProductID   uint64
	Name        string
	BrandName   string
	Description string
	Price       uint64
}

type UpdateProductDetailHandler CommandHandler[UpdateProductDetail]

type updateProductDetailHandler struct {
	logger      logger.Logger
	productRepo domain.ProductRepository
}

func NewUpdateProductDetailHandler(logger logger.Logger, productRepo domain.ProductRepository) UpdateProductDetailHandler {
	return &updateProductDetailHandler{
		logger:      logger,
		productRepo: productRepo,
	}
}

func (h *updateProductDetailHandler) Handle(ctx context.Context, cmd UpdateProductDetail) error {
	err := h.productRepo.UpdateProductDetail(ctx, cmd.ProductID, &valueobject.ProductDetail{
		Name:        cmd.Name,
		Description: cmd.Description,
		BrandName:   cmd.BrandName,
		Price:       cmd.Price,
	})

	if err != nil {
		return err
	}

	return nil
}
