package models

type Cart struct {
	Items         []CartItem `json:"items"`
	TotalPrice    float64    `json:"total_price"`
	TotalDiscount float64    `json:"total_discount"`
	FinalPrice    float64    `json:"final_price"`
}

type CartItem struct {
	ProductID     string  `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
	TotalDiscount float64 `json:"total_discount"`
}
