package repository

import (
	"context"

	"gorm.io/gorm"
)

type (
	Query[T any]    func(ctx context.Context, gormDB *gorm.DB) ([]T, error)
	QueryOne[T any] func(ctx context.Context, gormDB *gorm.DB) (*T, error)
)

type Repository[T any] interface {
	Execute(ctx context.Context, query Query[T]) ([]T, error)
	ExecuteOne(ctx context.Context, query QueryOne[T]) (*T, error)
}

type repository[T any] struct {
	gormDB *gorm.DB
}

func New[T any](gormDB *gorm.DB) Repository[T] {
	return &repository[T]{gormDB: gormDB}
}

func (repository *repository[T]) Execute(ctx context.Context, query Query[T]) ([]T, error) {
	result, err := query(ctx, repository.gormDB)
	return result, err
}

func (repository *repository[T]) ExecuteOne(ctx context.Context, query QueryOne[T]) (*T, error) {
	result, err := query(ctx, repository.gormDB)
	return result, err
}
