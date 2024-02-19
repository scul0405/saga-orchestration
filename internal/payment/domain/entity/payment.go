package entity

type Payment struct {
	ID           uint64
	CustomerID   uint64
	CurrencyCode string
	Amount       uint64
}
