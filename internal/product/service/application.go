package service

import (
	"github.com/scul0405/saga-orchestration/internal/product/app"
	"github.com/scul0405/saga-orchestration/internal/product/app/command"
	"github.com/scul0405/saga-orchestration/internal/product/app/query"
	"github.com/scul0405/saga-orchestration/internal/product/domain"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
)

func NewProductService(sf sonyflake.IDGenerator, logger logger.Logger, productRepo domain.ProductRepository) app.ProductApplication {
	return app.ProductApplication{
		Commands: app.ProductCommands{
			CreateProduct:            command.NewCreateProductHandler(sf, logger, productRepo),
			UpdateProductDetail:      command.NewUpdateProductDetailHandler(logger, productRepo),
			UpdateProductInventory:   command.NewUpdateProductInventoryHandler(logger, productRepo),
			RollbackProductInventory: command.NewRollbackProductInventoryHandler(logger, productRepo),
		},
		Queries: app.ProductQueries{
			CheckProducts: query.NewCheckProductsHandler(logger, productRepo),
			GetProducts:   query.NewGetProductsHandler(logger, productRepo),
			ListProducts:  query.NewListProductsHandler(logger, productRepo),
		},
	}
}

func NewCategoryService(sf sonyflake.IDGenerator, logger logger.Logger, categoryRepo domain.CategoryRepository) app.CategoryApplication {
	return app.CategoryApplication{
		Commands: app.CategoryCommands{
			CreateCategory: command.NewCreateCategoryHandler(sf, logger, categoryRepo),
		},
	}
}