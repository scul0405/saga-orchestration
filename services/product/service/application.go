package service

import (
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/services/product/app"
	command2 "github.com/scul0405/saga-orchestration/services/product/app/command"
	query2 "github.com/scul0405/saga-orchestration/services/product/app/query"
	"github.com/scul0405/saga-orchestration/services/product/domain"
)

func NewProductService(sf sonyflake.IDGenerator, logger logger.Logger, productRepo domain.ProductRepository) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreateProduct:            command2.NewCreateProductHandler(sf, logger, productRepo),
			UpdateProductDetail:      command2.NewUpdateProductDetailHandler(logger, productRepo),
			UpdateProductInventory:   command2.NewUpdateProductInventoryHandler(logger, productRepo),
			RollbackProductInventory: command2.NewRollbackProductInventoryHandler(logger, productRepo),
		},
		Queries: app.Queries{
			CheckProducts: query2.NewCheckProductsHandler(logger, productRepo),
			GetProducts:   query2.NewGetProductsHandler(logger, productRepo),
			ListProducts:  query2.NewListProductsHandler(logger, productRepo),
		},
	}
}
