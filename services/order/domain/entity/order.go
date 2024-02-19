package entity

import (
	"github.com/scul0405/saga-orchestration/services/order/domain/valueobject"
)

type Order struct {
	ID                uint64
	CustomerID        uint64
	PurchasedProducts *[]valueobject.PurchasedProduct
}
