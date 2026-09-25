package order

import (
	orderitem "kitchen-api/internal/app/order_item"
	"kitchen-api/internal/app/restaurant"
	"kitchen-api/internal/app/table"
	"kitchen-api/internal/enums"
	"kitchen-api/internal/helper"

	"github.com/google/uuid"
)

type Order struct {
	ID           uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RestaurantID uuid.UUID               `gorm:"type:uuid;not null" json:"restaurantId"`
	Restaurant   *restaurant.Restaurants `json:"restaurant,omitempty"`
	TableID      uuid.UUID               `gorm:"type:uuid;not null" json:"tableId"`
	Table        *table.Table            `json:"table,omitempty"`
	Status          enums.OrderStatus      `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"`
	TotalAmount     float64                `gorm:"type:decimal(10,2);not null;default:0" json:"totalAmount"`
	Notes           string                 `gorm:"type:text" json:"notes,omitempty"`
	PaymentStatus   string                 `gorm:"type:varchar(50);default:'UNPAID';not null" json:"paymentStatus"`
	StripeSessionID *string                `gorm:"type:varchar(255)" json:"stripeSessionId,omitempty"`
	OrderItems      []orderitem.OrderItem  `gorm:"foreignKey:OrderID" json:"orderItems,omitempty"`
	helper.Generic
}
