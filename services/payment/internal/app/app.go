package app

import (
	"github.com/scul0405/saga-orchestration/services/payment/internal/app/command"
	"github.com/scul0405/saga-orchestration/services/payment/internal/app/query"
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
