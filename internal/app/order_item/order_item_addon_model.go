package orderitem

import (
	"kitchen-api/internal/app/addon"
	"kitchen-api/internal/helper"

	"github.com/google/uuid"
)

type OrderItemAddon struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrderItemID uuid.UUID    `gorm:"type:uuid;not null;index" json:"orderItemId"`
	AddonID     uuid.UUID    `gorm:"type:uuid;not null" json:"addonId"`
	Addon       *addon.Addon `json:"addon,omitempty"`
	Quantity    int          `gorm:"not null;default:1" json:"quantity"`
	UnitPrice   float64      `gorm:"type:decimal(10,2);not null" json:"unitPrice"`
	helper.Generic
}
