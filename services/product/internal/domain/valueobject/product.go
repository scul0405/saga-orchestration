package valueobject

// ProductDetail value object
type ProductDetail struct {
	Name        string
	Description string
	BrandName   string
	Price       uint64
}

// PurchasedProduct value object
type PurchasedProduct struct {
	ID       uint64
	Quantity uint64
}
