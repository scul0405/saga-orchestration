package app

import (
	command2 "github.com/scul0405/saga-orchestration/services/order/app/command"
	"github.com/scul0405/saga-orchestration/services/order/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateOrder command2.CreateOrderHandler
	DeleteOrder command2.DeleteOrderHandler
}

type Queries struct {
	GetDetailedOrder query.GetDetailedOrderHandler
}
