package payment

import (
	"kitchen-api/internal/app/order"
)

type CreateCheckoutResponse struct {
	CheckoutURL string  `json:"checkoutUrl"`
	SessionID   string  `json:"sessionId"`
	TotalAmount float64 `json:"totalAmount"`
}

type SessionStatusResponse struct {
	SessionID string       `json:"sessionId"`
	Status    string       `json:"status"` // PENDING, COMPLETED, EXPIRED
	OrderID   *string      `json:"orderId,omitempty"`
	Order     *order.Order `json:"order,omitempty"`
}
