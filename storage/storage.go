package storage

import "context"

type Storage interface {
	Save(ctx context.Context, content []byte, key string) error
	Get(ctx context.Context, key string) ([]byte, error)
	Delete(ctx context.Context, key string) error
}
