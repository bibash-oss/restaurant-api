package payment

import (
	"github.com/google/uuid"
	"kitchen-api/internal/app/order"
	"kitchen-api/internal/app/restaurant"
	"kitchen-api/internal/app/table"
	"kitchen-api/internal/helper"
)

type PaymentSession struct {
	ID              uuid.UUID               `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StripeSessionID string                  `gorm:"type:varchar(255);uniqueIndex;not null" json:"stripeSessionId"`
	RestaurantID    uuid.UUID               `gorm:"type:uuid;not null" json:"restaurantId"`
	Restaurant      *restaurant.Restaurants `json:"restaurant,omitempty"`
	TableID         uuid.UUID               `gorm:"type:uuid;not null" json:"tableId"`
	Table           *table.Table            `json:"table,omitempty"`
	OrderPayload    string                  `gorm:"type:text;not null" json:"-"`
	TotalAmount     float64                 `gorm:"type:decimal(10,2);not null" json:"totalAmount"`
	Currency        string                  `gorm:"type:varchar(10);default:'aud';not null" json:"currency"`
	Status          string                  `gorm:"type:varchar(50);default:'PENDING';not null" json:"status"` // PENDING, COMPLETED, EXPIRED
	OrderID         *uuid.UUID              `gorm:"type:uuid" json:"orderId,omitempty"`
	Order           *order.Order            `json:"order,omitempty"`
	helper.Generic
}
