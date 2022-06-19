package repository

import (
	"gorm.io/gorm"
)

type (
	Query[T any]    func(gormDB *gorm.DB) ([]T, error)
	QueryOne[T any] func(gormDB *gorm.DB) (T, error)
)

type Repository[T any] interface {
	Execute(query Query[T]) ([]T, error)
	ExecuteOne(query QueryOne[T]) (T, error)
}

type repository[T any] struct {
	gormDB *gorm.DB
}

func New[T any](gormDB *gorm.DB) Repository[T] {
	return &repository[T]{gormDB: gormDB}
}

func (repository *repository[T]) Execute(query Query[T]) ([]T, error) {
	result, err := query(repository.gormDB)
	return result, err
}

func (repository *repository[T]) ExecuteOne(query QueryOne[T]) (T, error) {
	result, err := query(repository.gormDB)
	return result, err
}
