package entity

type Order struct {
	CustomerID uint64
	OrderItems *[]OrderItem
}
