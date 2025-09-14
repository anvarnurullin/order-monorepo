package model

import "time"

type Product struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	SKU          string    `json:"sku"`
	Price        float64   `json:"price"`
	QtyAvailable int       `json:"qty_available"`
	CreatedAt    time.Time `json:"created_at"`
}
