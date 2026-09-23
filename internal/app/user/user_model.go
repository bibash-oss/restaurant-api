package user

import (
	"kitchen-api/internal/app/restaurant"
	"kitchen-api/internal/enums"
	"kitchen-api/internal/helper"

	"github.com/google/uuid"
)

type Users struct {
	ID           uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string                 `gorm:"type:varchar(255);not null" json:"name"`
	Email        string                 `gorm:"type:varchar(255);not null;unique" json:"email"`
	Password     string                 `gorm:"type:varchar(255);not null" json:"-"`
	IsActive     bool                   `gorm:"type:boolean;default:false;not null" json:"isActive"`
	Role         enums.UserRole         `gorm:"type:user_role;not null;default:'USER'" json:"role"`
	RestaurantID *uuid.UUID              `gorm:"type:uuid;column:restaurant_id" json:"restaurantId,omitempty"`
	Restaurant   *restaurant.Restaurants `json:"restaurant,omitempty"`
	helper.Generic
}
