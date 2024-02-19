package service

import (
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/services/product/internal/app"
	"github.com/scul0405/saga-orchestration/services/product/internal/app/command"
	"github.com/scul0405/saga-orchestration/services/product/internal/app/query"
	"github.com/scul0405/saga-orchestration/services/product/internal/domain"
)

func NewProductService(sf sonyflake.IDGenerator, logger logger.Logger, productRepo domain.ProductRepository) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateProduct:            command.NewCreateProductHandler(sf, logger, productRepo),
			UpdateProductDetail:      command.NewUpdateProductDetailHandler(logger, productRepo),
			UpdateProductInventory:   command.NewUpdateProductInventoryHandler(logger, productRepo),
			RollbackProductInventory: command.NewRollbackProductInventoryHandler(logger, productRepo),
		},
		Queries: app.Queries{
			CheckProducts: query.NewCheckProductsHandler(logger, productRepo),
			GetProducts:   query.NewGetProductsHandler(logger, productRepo),
			ListProducts:  query.NewListProductsHandler(logger, productRepo),
		},
	}
}
