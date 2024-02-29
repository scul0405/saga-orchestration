package aggregate

import (
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/entity"
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/valueobject"
)

// Purchase aggregate
type Purchase struct {
	ID      uint64
	Order   *entity.Order
	Payment *valueobject.Payment
}
