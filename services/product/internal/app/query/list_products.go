package query

import (
	"context"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/services/product/internal/domain"
	"github.com/scul0405/saga-orchestration/services/product/internal/domain/valueobject"
)

type ListProducts struct {
	Limit  uint64
	Offset uint64
}

type ListProductsHandler QueryHandler[ListProducts, *[]valueobject.ProductCatalog]

type listProductsHandler struct {
	logger      logger.Logger
	productRepo domain.ProductRepository
}

func NewListProductsHandler(logger logger.Logger, productRepo domain.ProductRepository) ListProductsHandler {
	return &listProductsHandler{
		logger:      logger,
		productRepo: productRepo,
	}
}

func (h *listProductsHandler) Handle(ctx context.Context, query ListProducts) (*[]valueobject.ProductCatalog, error) {
	products, err := h.productRepo.ListProducts(ctx, query.Offset, query.Offset)
	if err != nil {
		return nil, err
	}

	return products, nil
}
