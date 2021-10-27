package storage

import (
	"context"

	"github.com/dotkom/image-server/models"
)

// Interface for storing images to a file system.
type FileStorage interface {
	// Save the given byte buffer with the given key
	Save(ctx context.Context, content []byte, key string) error

	// Retrieve a byte buffer with the given key
	Get(ctx context.Context, key string) ([]byte, error)

	// Delete the object at location of the given key
	Delete(ctx context.Context, key string) error
}

// Interface for storing images to a file system.
type MetaStorage interface {
	// Save the given image meta
	Save(ctx context.Context, meta models.ImageMeta) error

	// Retrieve a the given meta with the given key
	Get(ctx context.Context, key string) (*models.ImageMeta, error)

	// Delete the meta with the given key
	Delete(ctx context.Context, key string) error
}
