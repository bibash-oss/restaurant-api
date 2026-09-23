package orderitem

import (
	menuitem "kitchen-api/internal/app/menu_item"
	"kitchen-api/internal/helper"

	"github.com/google/uuid"
)

type OrderItem struct {
	ID         uuid.UUID         `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderID    uuid.UUID          `gorm:"type:uuid;not null" json:"orderId"`
	MenuItemID uuid.UUID          `gorm:"type:uuid;not null" json:"menuItemId"`
	MenuItem   *menuitem.MenuItem `json:"menuItem,omitempty"`
	Quantity   int               `gorm:"not null" json:"quantity"`
	UnitPrice  float64           `gorm:"type:decimal(10,2);not null" json:"unitPrice"`
	Addons     []OrderItemAddon  `gorm:"foreignKey:OrderItemID;constraint:OnDelete:CASCADE" json:"addons,omitempty"`
	helper.Generic
}

