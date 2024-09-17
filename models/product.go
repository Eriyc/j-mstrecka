package models

import "time"

type Product struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	TotalStock int    `json:"total_stock"`
}

type ProductPrice struct {
	ID            int64     `json:"id"`
	ProductID     int64     `json:"product_id"`
	PurchasePrice float64   `json:"purchase_price"`
	InternalPrice float64   `json:"internal_price"`
	ExternalPrice float64   `json:"external_price"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
}

type ProductWithPrice struct {
	Product Product
	Price   ProductPrice
}
