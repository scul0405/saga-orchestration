package app

import (
	command2 "github.com/scul0405/saga-orchestration/services/product/app/command"
	query2 "github.com/scul0405/saga-orchestration/services/product/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateProduct            command2.CreateProductHandler
	UpdateProductDetail      command2.UpdateProductDetailHandler
	UpdateProductInventory   command2.UpdateProductInventoryHandler
	RollbackProductInventory command2.RollbackProductInventoryHandler
}

type Queries struct {
	CheckProducts query2.CheckProductsHandler
	GetProducts   query2.GetProductsHandler
	ListProducts  query2.ListProductsHandler
}
