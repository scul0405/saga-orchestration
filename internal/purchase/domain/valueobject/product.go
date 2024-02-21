package valueobject

// Status enumeration
type Status int

const (
	// ProductOk is ok status
	ProductOk Status = iota
	// ProductNotFound is not found status
	ProductNotFound
)

type ProductStatus struct {
	ProductID uint64
	Status    Status
}
