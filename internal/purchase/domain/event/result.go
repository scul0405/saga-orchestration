package event

import "time"

var (
	StepUpdateProductInventory = "UPDATE_PRODUCT_INVENTORY"
	StepCreateOrder            = "CREATE_ORDER"
	StepCreatePayment          = "CREATE_PAYMENT"

	StatusExecute        = "EXUCUTE"
	StatusSucess         = "SUCCESS"
	StatusFailed         = "FAILED"
	StatusRollbacked     = "ROLLBACKED"
	StatusRollbackFailed = "ROLLBACK_FAIL"
)

// PurchaseResult event
type PurchaseResult struct {
	PurchaseID uint64
	Step       string
	Status     string
	Timestamp  time.Time
}
