package dto

type AddProductToCartRequest struct {
	CustomerID int    `json:"customer_id"`
	CartName   string `json:"cart_name"`
	ProductID  int    `json:"product_id"`
	Quantity   int    `json:"quantity"`
}
