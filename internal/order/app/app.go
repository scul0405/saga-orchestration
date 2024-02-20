package app

import (
	"github.com/scul0405/saga-orchestration/internal/order/app/command"
	"github.com/scul0405/saga-orchestration/internal/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command.CreateOrderHandler
	DeleteOrder command.DeleteOrderHandler
}

type Queries struct {
	GetDetailedOrder query.GetDetailedOrderHandler
}
