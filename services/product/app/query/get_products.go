package query

import (
	"context"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/services/product/domain"
	"github.com/scul0405/saga-orchestration/services/product/domain/entity"
)

type GetProducts struct {
	ProductIDs *[]uint64
}

type GetProductsHandler QueryHandler[GetProducts, *[]entity.Product]

type getProductsHandler struct {
	logger      logger.Logger
	productRepo domain.ProductRepository
}

func NewGetProductsHandler(logger logger.Logger, productRepo domain.ProductRepository) GetProductsHandler {
	return &getProductsHandler{
		logger:      logger,
		productRepo: productRepo,
	}
}

func (h *getProductsHandler) Handle(ctx context.Context, query GetProducts) (*[]entity.Product, error) {
	products := make([]entity.Product, len(*(query.ProductIDs)))

	for i, productID := range *(query.ProductIDs) {
		product, err := h.productRepo.GetProduct(ctx, productID)
		if err != nil {
			return nil, err
		}
		products[i] = *product
	}

	return &products, nil
}
