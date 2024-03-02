package eventhandler

import (
	"github.com/scul0405/saga-orchestration/internal/payment/app/command"
	pb "github.com/scul0405/saga-orchestration/proto"
)

func decodePb2CreatePaymentCmd(purchase *pb.CreatePurchaseRequest) command.CreatePayment {
	return command.CreatePayment{
		ID:           purchase.PurchaseId,
		CustomerID:   purchase.Purchase.Order.CustomerId,
		Amount:       purchase.Purchase.Payment.Amount,
		CurrencyCode: purchase.Purchase.Payment.CurrencyCode,
	}
}
