package query

import (
	"context"
	"github.com/scul0405/saga-orchestration/internal/order/domain"
	"github.com/scul0405/saga-orchestration/internal/order/domain/valueobject"
	"github.com/scul0405/saga-orchestration/internal/order/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/pkg/logger"
)

type GetDetailedOrder struct {
	OrderID uint64
}

type GetDetailedOrderHandler QueryHandler[GetDetailedOrder, *valueobject.DetailedOrder]

type getDetailedOrderHandler struct {
	logger     logger.Logger
	orderRepo  domain.OrderRepository
	productSvc grpc.ProductService
}

func NewGetDetailedOrderHandler(logger logger.Logger, orderRepo domain.OrderRepository, productSvc grpc.ProductService) GetDetailedOrderHandler {
	return &getDetailedOrderHandler{
		logger:     logger,
		orderRepo:  orderRepo,
		productSvc: productSvc,
	}
}

func (h *getDetailedOrderHandler) Handle(ctx context.Context, query GetDetailedOrder) (*valueobject.DetailedOrder, error) {
	order, err := h.orderRepo.GetOrder(ctx, query.OrderID)
	if err != nil {
		return nil, err
	}

	productIds := make([]uint64, len(*order.PurchasedProducts))
	for i, p := range *order.PurchasedProducts {
		productIds[i] = p.ID
	}

	products, err := h.productSvc.GetProducts(ctx, &productIds)
	if err != nil {
		return nil, err
	}

	return &valueobject.DetailedOrder{
		ID:                order.ID,
		CustomerID:        order.CustomerID,
		PurchasedProducts: products,
	}, nil
}
