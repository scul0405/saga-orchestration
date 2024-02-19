package dto

type Product struct {
	ID          uint64 `json:"id"`
	CategoryID  uint64 `json:"category_id"`
	Name        string `json:"name"`
	BrandName   string `json:"brand_name"`
	Description string `json:"description"`
	Price       uint64 `json:"price"`
	Quantity    uint64 `json:"quantity"`
}
