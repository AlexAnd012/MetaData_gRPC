package storage

import (
	"context"
	mediametav1 "github.com/AlexAnd012/mediameta/gen/go/mediameta/v1"
)

// Repository Чтобы gRPC-методы не зависели от Postgres работаем с абстракцией Repository
type Repository interface {
	Insert(ctx context.Context, m *mediametav1.Metadata) error
	Get(ctx context.Context, id string) (*mediametav1.Metadata, error)
	List(ctx context.Context, limit, offset int) ([]*mediametav1.Metadata, int, error)
	Update(ctx context.Context, m *mediametav1.Metadata) error
	Delete(ctx context.Context, id string) error // hard delete
}
