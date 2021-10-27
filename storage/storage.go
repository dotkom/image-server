package storage

import "context"

// Interface for storing images to a file system.
type Storage interface {
	// Save the given byte buffer with the given key
	Save(ctx context.Context, content []byte, key string) error

	// Retrieve a byte buffer with the given key
	Get(ctx context.Context, key string) ([]byte, error)

	// Delete the object at location of the given key
	Delete(ctx context.Context, key string) error
}
