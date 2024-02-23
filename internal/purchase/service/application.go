package service

import (
	"github.com/scul0405/saga-orchestration/internal/purchase/app"
	"github.com/scul0405/saga-orchestration/internal/purchase/app/command"
	"github.com/scul0405/saga-orchestration/internal/purchase/app/query"
	"github.com/scul0405/saga-orchestration/internal/purchase/eventhandler"
	"github.com/scul0405/saga-orchestration/internal/purchase/infrastructure/grpc"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
)

func NewPurchaseService(sf sonyflake.IDGenerator, logger logger.Logger, productSvc grpc.ProductService, evPub eventhandler.PurchaseEventHandler) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreatePurchase: command.NewCreatePurchaseHandler(sf, logger, productSvc, evPub),
		},
		Queries: app.Queries{
			CheckProducts: query.NewCheckProductsHandler(logger, productSvc),
		},
	}
}
