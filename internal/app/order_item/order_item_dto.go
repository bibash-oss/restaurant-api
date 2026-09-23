package orderitem

type CreateOrderItemInput struct {
	MenuItemID string `json:"menuItemId" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,gt=0"`
}

type OrderItemResponse struct {
	ID         string  `json:"id"`
	OrderID    string  `json:"orderId"`
	MenuItemID string  `json:"menuItemId"`
	ItemName   string  `json:"itemName,omitempty"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unitPrice"`
	Subtotal   float64 `json:"subtotal"`
}
