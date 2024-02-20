package dto

type Payment struct {
	ID           uint64 `json:"id"`
	CustomerID   uint64 `json:"customer_id"`
	Amount       uint64 `json:"amount"`
	CurrencyCode string `json:"currency_code"`
}
