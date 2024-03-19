package app

import (
	"github.com/scul0405/saga-orchestration/internal/product/app/command"
	"github.com/scul0405/saga-orchestration/internal/product/app/query"
)

type ProductApplication struct {
	Commands ProductCommands
	Queries  ProductQueries
}

type CategoryApplication struct {
	Commands CategoryCommands
	Queries CategoryQueries
}

type ProductCommands struct {
	CreateProduct            command.CreateProductHandler
	UpdateProductDetail      command.UpdateProductDetailHandler
	UpdateProductInventory   command.UpdateProductInventoryHandler
	RollbackProductInventory command.RollbackProductInventoryHandler
}

type ProductQueries struct {
	CheckProducts query.CheckProductsHandler
	GetProducts   query.GetProductsHandler
	ListProducts  query.ListProductsHandler
}

type CategoryCommands struct {
	CreateCategory command.CreateCategoryHandler
}

type CategoryQueries struct {

}