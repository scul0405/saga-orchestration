package dto

type Purchase struct {
	OrderItems *[]OrderItem `json:"order_items"`
	Payment    *Payment     `json:"payment"`
}

type OrderItem struct {
	ProductID uint64 `json:"product_id"`
	Quantity  uint64 `json:"quantity"`
}

type Payment struct {
	CurrencyCode string `json:"currency_code"`
}
