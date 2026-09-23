package addon

import (
	"kitchen-api/internal/app/restaurant"
	"kitchen-api/internal/helper"

	"github.com/google/uuid"
)

type Addon struct {
	ID           uuid.UUID               `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RestaurantID uuid.UUID               `gorm:"type:uuid;not null" json:"restaurantId"`
	Restaurant   *restaurant.Restaurants `json:"restaurant,omitempty"`
	Name         string                  `gorm:"type:varchar(255);not null" json:"name"`
	Price        float64                 `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive     bool                    `gorm:"type:boolean;default:true;not null" json:"isActive"`
	helper.Generic
}

