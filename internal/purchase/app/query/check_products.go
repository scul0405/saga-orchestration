package query

import (
	"context"
	"errors"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/valueobject"
	"github.com/scul0405/saga-orchestration/internal/purchase/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/pkg/logger"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrProductNotEnough = errors.New("product not enough")
	ErrInternal         = errors.New("internal error")
)

type CheckProducts struct {
	OrderItems *[]OrderItem
}

type OrderItem struct {
	ID       uint64
	Quantity uint64
}

type CheckProductsHandler QueryHandler[CheckProducts, *[]valueobject.ProductStatus]

type checkProductsHandler struct {
	logger     logger.Logger
	productSvc grpc.ProductService
}

func NewCheckProductsHandler(logger logger.Logger, productSvc grpc.ProductService) CheckProductsHandler {
	return &checkProductsHandler{
		logger:     logger,
		productSvc: productSvc,
	}
}

func (h *checkProductsHandler) Handle(ctx context.Context, query CheckProducts) (*[]valueobject.ProductStatus, error) {
	orderItems := make([]entity.OrderItem, len(*query.OrderItems))
	for i, item := range *query.OrderItems {
		orderItems[i] = entity.OrderItem{
			ID:       item.ID,
			Quantity: item.Quantity,
		}
	}

	productStatuses, err := h.productSvc.CheckProducts(ctx, &orderItems)
	if err != nil {
		return nil, err
	}

	for _, status := range *productStatuses {
		switch status.Status {
		case valueobject.ProductNotFound:
			return nil, ErrProductNotFound
		case valueobject.ProductNotEnough:
			return nil, ErrProductNotEnough
		case valueobject.ProductInternalError:
			return nil, ErrInternal
		default:
			continue
		}
	}

	return productStatuses, nil
}
