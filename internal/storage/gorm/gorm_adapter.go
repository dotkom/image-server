package gorm_adapter

import (
	"context"
	"time"

	"github.com/dotkom/image-server/internal/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormAdapter struct {
	db *gorm.DB
}

type DBDriver string

const (
	PostgresDriver DBDriver = "postgres"
	SqliteDriver   DBDriver = "sqlite"
)

type imageMeta struct {
	Key         uuid.UUID `gorm:"primaryKey;type:uuid;"`
	Name        string
	Description string
	Tags        []string `gorm:"type:text[]"`
	Mime        string
	Size        uint64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func New(driver DBDriver, dsn string) *GormAdapter {
	log.Info("Creating GORM MetaStorage adapter")
	adapter := &GormAdapter{}
	var err error
	switch driver {
	case PostgresDriver:
		adapter.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case SqliteDriver:
		adapter.db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}
	if err != nil {
		log.Fatal("Failed to establish database connection", err)
	}
	return adapter
}

// Migrate the database to the current schema as defined by the models in this module
func (adapter *GormAdapter) Migrate() {
	log.Info("Migrating database models")
	adapter.db.AutoMigrate(&imageMeta{})
}

func (adapter *GormAdapter) Save(ctx context.Context, meta models.ImageMeta) error {
	var model imageMeta

	adapter.db.Model(&model).Create(&imageMeta{
		Key:         meta.Key,
		Name:        meta.Name,
		Description: meta.Description,
		Tags:        meta.Tags,
		Mime:        meta.Mime,
		Size:        meta.Size,
	})
	return nil
}

func (adapter *GormAdapter) Get(ctx context.Context, key string) (*models.ImageMeta, error) {
	var meta imageMeta
	result := adapter.db.Limit(1).Find(&meta, "key = ?", key)
	if result.Error != nil {
		return nil, result.Error
	}
	return &models.ImageMeta{
		Key:         meta.Key,
		Name:        meta.Name,
		Description: meta.Description,
		Tags:        meta.Tags,
		Mime:        meta.Mime,
		Size:        meta.Size,
	}, nil
}

func (adapter *GormAdapter) Delete(ctx context.Context, key string) error {
	var meta imageMeta
	result := adapter.db.Delete(&meta, "key = ?", key)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
