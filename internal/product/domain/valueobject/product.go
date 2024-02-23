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

// ProductStatus value object
type ProductStatus struct {
	ID     uint64
	Price  uint64
	Status bool
}

// ProductCatalog value object
type ProductCatalog struct {
	ID         uint64
	CategoryID uint64
	Name       string
	Price      uint64
	Inventory  uint64
}
