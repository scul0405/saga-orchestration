package service

import (
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/services/order/app"
	command2 "github.com/scul0405/saga-orchestration/services/order/app/command"
	"github.com/scul0405/saga-orchestration/services/order/app/query"
	"github.com/scul0405/saga-orchestration/services/order/domain"
	"github.com/scul0405/saga-orchestration/services/order/infrastructure/grpc/product"
)

func NewOrderService(sf sonyflake.IDGenerator, logger logger.Logger, orderRepo domain.OrderRepository, productSvc product.ProductService) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command2.NewCreateOrderHandler(sf, logger, orderRepo),
			DeleteOrder: command2.NewDeleteOrderHandler(logger, orderRepo),
		},
		Queries: app.Queries{
			GetDetailedOrder: query.NewGetDetailedOrderHandler(logger, orderRepo, productSvc),
		},
	}
}
