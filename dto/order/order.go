package orderdto

type CreateOrderRequest struct {
	ProductID int ` json:"product_id" `
	CartID    int ` json:"cart_id"`
	Qty       int ` json:"qty" `
	SubAmount int ` json:"sub_amount"`
}

type OrderUpdate struct {
	Qty       int ` json:"qty" `
	SubAmount int ` json:"sub_amount"`
}
