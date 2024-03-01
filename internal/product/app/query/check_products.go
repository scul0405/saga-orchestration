package query

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/product/domain"
	"github.com/scul0405/saga-orchestration/internal/product/domain/valueobject"
	"github.com/scul0405/saga-orchestration/pkg/logger"
)

type CheckProducts struct {
	Items *[]CheckItem
}

type CheckItem struct {
	ProductID uint64
	Quantity  uint64
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
	productStatuses := make([]valueobject.ProductStatus, len(*(query.Items)))

	for i, item := range *(query.Items) {
		status, err := h.productRepo.CheckProduct(ctx, item.ProductID, item.Quantity)
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
