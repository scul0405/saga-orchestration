package eventhandler

import (
	"github.com/scul0405/saga-orchestration/internal/order/app/command"
	pb "github.com/scul0405/saga-orchestration/proto"
)

func decodePb2CreateOrderCmd(purchase *pb.CreatePurchaseRequest) command.CreateOrder {
	purchasedProduct := make([]command.PurchasedProduct, len(purchase.Purchase.Order.OrderItems))
	for i, item := range purchase.Purchase.Order.OrderItems {
		purchasedProduct[i] = command.PurchasedProduct{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		}
	}

	return command.CreateOrder{
		OrderID:    purchase.PurchaseId,
		CustomerID: purchase.Purchase.Order.CustomerId,
		Products:   &purchasedProduct,
	}
}
