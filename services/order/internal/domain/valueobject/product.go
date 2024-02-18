package valueobject

// PurchasedProduct value object
type PurchasedProduct struct {
	ID       uint64
	Quantity uint64
}

// DetailedPurchasedProduct value object
type DetailedPurchasedProduct struct {
	ID          uint64
	CategoryID  uint64
	Name        string
	BrandName   string
	Description string
	Price       uint64
	Quantity    uint64
}
