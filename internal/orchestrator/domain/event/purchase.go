package event

import (
	"github.com/scul0405/saga-orchestration/internal/orchestrator/domain/aggregate"
)

type CreatePurchaseResponse struct {
	Purchase *aggregate.Purchase
	Success  bool
	Error    string
}

type RollbackResponse struct {
	CustomerID uint64
	PurchaseID uint64
	Success    bool
	Error      string
}
