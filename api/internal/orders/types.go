package orders

type CreateOrderParams struct {
	UserID int64
	items  []OrderItem
}

type OrderItem struct {
	ProductId int64
	Price     int64
	Quantity  int32
}
