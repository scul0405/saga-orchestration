package dto

type Product struct {
	ID          uint64 `json:"id"`
	CategoryID  uint64 `json:"category_id"`
	Name        string `json:"name"`
	BrandName   string `json:"brand_name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Inventory   uint64 `json:"inventory"`
}

type CreateProduct struct {
	CategoryID  uint64 `json:"category_id"`
	Name        string `json:"name"`
	BrandName   string `json:"brand_name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Inventory   uint64 `json:"inventory"`
}

type UpdateProductDetail struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	BrandName   string `json:"brand_name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
}
