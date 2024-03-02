package valueobject

// Status enumeration
type Status int

const (
	// ProductOk is ok status
	ProductOk Status = iota
	// ProductNotFound is not found status
	ProductNotFound
	// ProductNotEnough is not enough status
	ProductNotEnough
	// ProductInternalError is internal error status
	ProductInternalError
)

type ProductStatus struct {
	ProductID uint64
	Price     uint64
	Status    Status
}
