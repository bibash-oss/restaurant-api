package table

import (
	"kitchen-api/internal/app/restaurant"
	"kitchen-api/internal/helper"

	"github.com/google/uuid"
)

type Table struct {
	ID           uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RestaurantID uuid.UUID               `gorm:"type:uuid;not null" json:"restaurantId"`
	Restaurant   *restaurant.Restaurants `json:"restaurant,omitempty"`
	Number       string                 `gorm:"type:varchar(50);not null" json:"number"`
	IsActive     bool                   `gorm:"type:boolean;default:true;not null" json:"isActive"`
	helper.Generic
}

