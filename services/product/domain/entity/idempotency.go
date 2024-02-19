package entity

// Idempotency entity
type Idempotency struct {
	ID        uint64 // id of a purchase
	ProductID uint64
	Quantity  uint64
}
