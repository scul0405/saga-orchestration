package model

type Product struct {
	ID          uint64 `gorm:"primaryKey"`
	CategoryID  uint64
	Category    Category `gorm:"foreignKey:CategoryID"`
	Name        string   `gorm:"type:varchar(256);not null"`
	Description string   `gorm:"type:text;not null"`
	BrandName   string   `gorm:"type:varchar(256);not null"`
	Inventory   uint64   `gorm:"not null"`
	Price       uint64   `gorm:"not null"`
	UpdatedAt   int64    `gorm:"autoUpdateTime:milli"`
	CreatedAt   int64    `gorm:"autoCreateTime:milli"`
}

type Idempotency struct {
	ID         uint64 `gorm:"primaryKey"`
	ProductID  uint64 `gorm:"primaryKey"`
	Quantity   uint64 `gorm:"not null"`
	Rollbacked bool   `gorm:"not null"`
	CreatedAt  int64  `gorm:"autoCreateTime:milli"`
}
