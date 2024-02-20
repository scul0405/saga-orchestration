package entity

import (
	"github.com/scul0405/saga-orchestration/internal/order/domain/valueobject"
)

type Order struct {
	ID                uint64
	CustomerID        uint64
	PurchasedProducts *[]valueobject.PurchasedProduct
}
