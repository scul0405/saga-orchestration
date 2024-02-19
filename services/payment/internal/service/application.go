package service

import (
	"github.com/scul0405/saga-orchestration/pkg/logger"
	"github.com/scul0405/saga-orchestration/pkg/sonyflake"
	"github.com/scul0405/saga-orchestration/services/payment/internal/app"
	"github.com/scul0405/saga-orchestration/services/payment/internal/app/command"
	"github.com/scul0405/saga-orchestration/services/payment/internal/app/query"
	"github.com/scul0405/saga-orchestration/services/payment/internal/domain"
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
