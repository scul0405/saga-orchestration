package eventhandler

import (
	"github.com/scul0405/saga-orchestration/internal/product/app/command"
	pb "github.com/scul0405/saga-orchestration/proto"
)

func decodePb2UpdateProductInventoryCmd(purchase *pb.CreatePurchaseRequest) command.UpdateProductInventory {
	purchasedProduct := make([]command.PurchasedProduct, len(purchase.Purchase.Order.OrderItems))
	for i, item := range purchase.Purchase.Order.OrderItems {
		purchasedProduct[i] = command.PurchasedProduct{
			ID:       item.ProductId,
			Quantity: item.Quantity,
		}
	}

	return command.UpdateProductInventory{
		IdempotencyKey:    purchase.PurchaseId,
		PurchasedProducts: &purchasedProduct,
	}
}

func decodePb2RollbackProductInventoryCmd(purchase *pb.CreatePurchaseRequest) command.RollbackProductInventory {
	purchasedProduct := make([]command.PurchasedProduct, len(purchase.Purchase.Order.OrderItems))
	for i, item := range purchase.Purchase.Order.OrderItems {
		purchasedProduct[i] = command.PurchasedProduct{
			ID:       item.ProductId,
			Quantity: item.Quantity,
		}
	}

	return command.RollbackProductInventory{
		IdempotencyKey:    purchase.PurchaseId,
		PurchasedProducts: &purchasedProduct,
	}
}
