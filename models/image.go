package models

import "github.com/google/uuid"

// Gorm model for storing information about images
type ImageMeta struct {
	Key         uuid.UUID
	Name        string
	Description string
	Tags        []string
	Mime        string
	Size        uint64
}
