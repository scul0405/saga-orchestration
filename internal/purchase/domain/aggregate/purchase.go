package aggregate

import (
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/purchase/domain/valueobject"
)

// Purchase aggregate
type Purchase struct {
	ID      uint64
	Order   *entity.Order
	Payment *valueobject.Payment
}
