package valueobject

// DetailedOrder value object
type DetailedOrder struct {
	ID                uint64
	CustomerID        uint64
	PurchasedProducts *[]DetailedPurchasedProduct
}
