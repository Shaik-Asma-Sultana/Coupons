package models

import "time"

type Coupon struct {
	ID      string        `json:"id"`
	Type    string        `json:"type"`
	Details CouponDetails `json:"details"`
}

type CouponDetails struct {
	Threshold        float64      `json:"threshold,omitempty"`
	Discount         float64      `json:"discount"`
	ProductID        string       `json:"product_id,omitempty"`
	BuyProducts      []BuyProduct `json:"buy_products,omitempty"`
	GetProducts      []GetProduct `json:"get_products,omitempty"`
	RepetitionLimit  int          `json:"repetition_limit,omitempty"`
	ExpiryDate       *time.Time   `json:"expiry_date,omitempty"`
	MaxUses          int          `json:"max_uses,omitempty"`
	Uses             int          `json:"uses,omitempty"`
	Exclusive        bool         `json:"exclusive,omitempty"`
	MinCartValue     float64      `json:"min_cart_value,omitempty"`
	ExcludedProducts []string     `json:"excluded_products,omitempty"`
}

type BuyProduct struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type GetProduct struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
