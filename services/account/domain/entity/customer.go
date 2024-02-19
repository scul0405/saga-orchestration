package entity

import (
	"github.com/scul0405/saga-orchestration/services/account/domain/valueobject"
)

// Customer entity
type Customer struct {
	ID           uint64
	Active       bool
	Password     string
	PersonalInfo *valueobject.CustomerPersonalInfo
	DeliveryInfo *valueobject.CustomerDeliveryInfo
}
