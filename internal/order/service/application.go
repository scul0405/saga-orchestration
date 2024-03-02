package service

import (
	"github.com/scul0405/saga-orchestration/internal/order/app"
	"github.com/scul0405/saga-orchestration/internal/order/app/command"
	"github.com/scul0405/saga-orchestration/internal/order/app/query"
	"github.com/scul0405/saga-orchestration/internal/order/domain"
	"github.com/scul0405/saga-orchestration/internal/order/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/pkg/logger"
)

func NewOrderService(logger logger.Logger, orderRepo domain.OrderRepository, productSvc grpc.ProductService) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateOrder: command.NewCreateOrderHandler(logger, orderRepo),
			DeleteOrder: command.NewDeleteOrderHandler(logger, orderRepo),
		},
		Queries: app.Queries{
			GetDetailedOrder: query.NewGetDetailedOrderHandler(logger, orderRepo, productSvc),
		},
	}
}
