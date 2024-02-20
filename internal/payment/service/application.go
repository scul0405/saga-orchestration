package service

import (
	"github.com/scul0405/saga-orchestration/internal/payment/app"
	"github.com/scul0405/saga-orchestration/internal/payment/app/command"
	"github.com/scul0405/saga-orchestration/internal/payment/app/query"
	"github.com/scul0405/saga-orchestration/internal/payment/domain"
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
)

func NewPaymentService(sf sonyflake.IDGenerator, logger logger.Logger, paymentRepo domain.PaymentRepository) app.Application {
	return app.Application{
		Commands: app.Commands{
			CreatePayment:   command.NewCreatePaymentHandler(sf, logger, paymentRepo),
			RollbackPayment: command.NewRollbackPaymentHandler(logger, paymentRepo),
		},
		Queries: app.Queries{
			GetPayment: query.NewGetPaymentHandler(logger, paymentRepo),
		},
	}
}
