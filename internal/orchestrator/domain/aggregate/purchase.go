package aggregate

import (
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/entity"
)

// Purchase aggregate
type Purchase struct {
	ID      uint64
	Order   *entity.Order
	Payment *entity.Payment
}
