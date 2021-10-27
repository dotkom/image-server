package models

import "gorm.io/gorm"

// Migrate the database to the current schema as defined by the models in this module
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&Image{})
}
