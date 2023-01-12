package models

type Product struct {
	PID         uint    `gorm:"primary_key"`
	ProductName string  `gorm:"unique; not null"`
	Description string  `gorm:"not null"`
	Price       float64 `gorm:"not null"`
	StockIn     bool    `gorm:"not null"`
}
