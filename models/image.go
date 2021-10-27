package models

import (
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Gorm model for storing information about images
type Image struct {
	gorm.Model
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;"`
	Name        string
	Description string
	Tags        pq.StringArray `gorm:"type:text[]"`
	Mime        string
	Size        uint64
}
