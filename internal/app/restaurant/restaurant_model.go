package restaurant

import (
	"kitchen-api/internal/helper"

	"github.com/google/uuid"
)

type Restaurants struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name     string    `gorm:"type:varchar(255);not null;unique" json:"name"`
	Slug     string    `gorm:"type:varchar(355);not null;unique" json:"slug"`
	IsActive bool      `gorm:"type:boolean;default:true;not null" json:"isActive"`
	Address  *string   `gorm:"type:varchar(255)" json:"address"`
	helper.Generic
}
