package orders

type CreateOrderParams struct {
	UserID int64 `json:"user_id"`
	Items  []OrderItem `json:"items"`
}

type OrderItem struct {
	ProductId int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}
