package menucategory

import (
	"kitchen-api/internal/app/restaurant"
	"kitchen-api/internal/helper"

	"github.com/google/uuid"
)

type MenuCategory struct {
	ID           uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RestaurantID uuid.UUID               `gorm:"type:uuid;not null" json:"restaurantId"`
	Restaurant   *restaurant.Restaurants `json:"restaurant,omitempty"`
	Name         string                 `gorm:"type:varchar(255);not null" json:"name"`
	helper.Generic
}
