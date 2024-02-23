package query

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/product/domain"
	"github.com/scul0405/saga-orchestration/internal/product/domain/valueobject"
	"github.com/scul0405/saga-orchestration/pkg/logger"
)

type CheckProducts struct {
	ProductIDs *[]uint64
}

type CheckProductsHandler QueryHandler[CheckProducts, *[]valueobject.ProductStatus]

type checkProductsHandler struct {
	logger      logger.Logger
	productRepo domain.ProductRepository
}

func NewCheckProductsHandler(logger logger.Logger, productRepo domain.ProductRepository) CheckProductsHandler {
	return &checkProductsHandler{
		logger:      logger,
		productRepo: productRepo,
	}
}

func (h *checkProductsHandler) Handle(ctx context.Context, query CheckProducts) (*[]valueobject.ProductStatus, error) {
	productStatuses := make([]valueobject.ProductStatus, len(*(query.ProductIDs)))

	for i, productID := range *(query.ProductIDs) {
		status, err := h.productRepo.CheckProduct(ctx, productID)
		if err != nil {
			return nil, err
		}
		productStatuses[i] = valueobject.ProductStatus{
			ID:     status.ID,
			Status: status.Status,
			Price:  status.Price,
		}
	}

	return &productStatuses, nil
}
