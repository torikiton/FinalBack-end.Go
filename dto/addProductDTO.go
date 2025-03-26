package dto

type AddProductToCartRequest struct {
	CustomerID int    `json:"customer_id"`
	CartName   string `json:"cart_name"`
	ProductID  int    `json:"product_id"`
	Quantity   int    `json:"quantity"`
}

type CartItemDTO struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"total_price"`
}

type CartDTO struct {
	CartName string        `json:"cart_name"`
	Items    []CartItemDTO `json:"items"`
}
type SearchProductRequest struct {
	MinPrice    float64 `json:"MinPrice"`
	MaxPrice    float64 `json:"MaxPrice"`
	ProductName string  `json:"ProductName"`
}
