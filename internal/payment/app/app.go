package app

import (
	"github.com/scul0405/saga-orchestration/internal/payment/app/command"
	"github.com/scul0405/saga-orchestration/internal/payment/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreatePayment   command.CreatePaymentHandler
	RollbackPayment command.RollbackPaymentHandler
}

type Queries struct {
	GetPayment query.GetPaymentHandler
}
