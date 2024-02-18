package service

import (
	"github.com/scul0405/saga-orchestration/services/order/internal/app"
	"github.com/scul0405/saga-orchestration/services/order/internal/app/command"
	"github.com/scul0405/saga-orchestration/services/order/internal/app/query"
	"github.com/scul0405/saga-orchestration/services/order/internal/domain"
	"github.com/scul0405/saga-orchestration/services/order/internal/infrastructure/grpc/product"
	"github.com/scul0405/saga-orchestration/services/order/internal/infrastructure/logger"
	"github.com/scul0405/saga-orchestration/services/order/pkg"
)

func NewOrderService(sf pkg.IDGenerator, logger logger.Logger, orderRepo domain.OrderRepository, productSvc product.ProductService) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(sf, logger, orderRepo),
			DeleteOrder: command.NewDeleteOrderHandler(logger, orderRepo),
		},
		Queries: app.Queries{
			GetDetailedOrder: query.NewGetDetailedOrderHandler(logger, orderRepo, productSvc),
		},
	}
}
