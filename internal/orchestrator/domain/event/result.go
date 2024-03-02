package event

import "time"

var (
	StepUpdateProductInventory = "UPDATE_PRODUCT_INVENTORY"
	StepCreateOrder            = "CREATE_ORDER"
	StepCreatePayment          = "CREATE_PAYMENT"

	StatusExecute        = "EXUCUTE"
	StatusSucess         = "SUCCESS"
	StatusFailed         = "FAILED"
	StatusRollback       = "ROLLBACK"
	StatusRollbackFailed = "ROLLBACK_FAILED"
)

// PurchaseResult event
type PurchaseResult struct {
	PurchaseID uint64
	Step       string
	Status     string
	Timestamp  time.Time
}
