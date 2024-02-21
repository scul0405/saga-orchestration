package entity

type Order struct {
	ID                uint64
	CustomerID        uint64
	PurchasedProducts *[]PurchasedProduct
}
