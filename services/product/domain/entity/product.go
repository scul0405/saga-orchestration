package entity

import (
	"github.com/scul0405/saga-orchestration/services/product/domain/valueobject"
)

// Product entity
type Product struct {
	ID         uint64
	CategoryID uint64
	Detail     *valueobject.ProductDetail
	Inventory  uint64
}
