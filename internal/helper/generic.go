package helper

import (
	"time"

	"github.com/google/uuid"
)

type Generic struct {
	CreatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:created_at" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP;column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"type:timestamp;column:deleted_at" json:"deletedAt,omitempty"`
	CreatedBy *uuid.UUID `gorm:"type:uuid;column:created_by" json:"createdBy,omitempty"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid;column:updated_by" json:"updatedBy,omitempty"`
	DeletedBy *uuid.UUID `gorm:"type:uuid;column:deleted_by" json:"deletedBy,omitempty"`
}
