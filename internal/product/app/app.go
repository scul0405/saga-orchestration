package app

import (
	"github.com/scul0405/saga-orchestration/internal/product/app/command"
	"github.com/scul0405/saga-orchestration/internal/product/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateProduct            command.CreateProductHandler
	UpdateProductDetail      command.UpdateProductDetailHandler
	UpdateProductInventory   command.UpdateProductInventoryHandler
	RollbackProductInventory command.RollbackProductInventoryHandler
}

type Queries struct {
	CheckProducts query.CheckProductsHandler
	GetProducts   query.GetProductsHandler
	ListProducts  query.ListProductsHandler
}
