package menuitem

import (
	"kitchen-api/internal/app/addon"
	menucategory "kitchen-api/internal/app/menu_category"
	"kitchen-api/internal/app/restaurant"
	"kitchen-api/internal/helper"

	"github.com/google/uuid"
)

type MenuItem struct {
	ID           uuid.UUID                 `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	RestaurantID uuid.UUID                  `gorm:"type:uuid;not null" json:"restaurantId"`
	Restaurant   *restaurant.Restaurants    `json:"restaurant,omitempty"`
	CategoryID   uuid.UUID                  `gorm:"type:uuid;not null" json:"categoryId"`
	Category     *menucategory.MenuCategory `json:"category,omitempty"`
	Name         string                    `gorm:"type:varchar(255);not null" json:"name"`
	Description  *string                   `gorm:"type:text" json:"description,omitempty"`
	Price        float64                   `gorm:"type:decimal(10,2);not null" json:"price"`
	ImageURL     *string                   `gorm:"type:varchar(500)" json:"imageUrl,omitempty"`
	IsActive     bool                      `gorm:"type:boolean;default:true;not null" json:"isActive"`
	Addons       []*addon.Addon            `gorm:"many2many:menu_item_addons;" json:"addons,omitempty"`
	helper.Generic
}
