package model

type Category struct {
	ID          uint64 `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(256);not null"`
	Description string `gorm:"type:text;not null"`
	UpdatedAt   int64  `gorm:"autoUpdateTime:milli"`
	CreatedAt   int64  `gorm:"autoCreateTime:milli"`
}
