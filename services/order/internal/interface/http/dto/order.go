package dto

type Order struct {
	OrderID  uint64    `json:"order_id"`
	Products []Product `json:"products"`
}
