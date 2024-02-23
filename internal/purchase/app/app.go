package app

import (
	"github.com/scul0405/saga-orchestration/internal/purchase/app/command"
	"github.com/scul0405/saga-orchestration/internal/purchase/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreatePurchase command.CreatePurchaseHandler
}

type Queries struct {
	CheckProducts query.CheckProductsHandler
}
