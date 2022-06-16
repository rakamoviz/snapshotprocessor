package repository

import (
	"gorm.io/gorm"
)

type (
	Query[T any] func(gormDB *gorm.DB) (*T, error)
)

type Repository[T any] interface {
	Execute(query Query[T]) (*T, error)
}

type repositoryStruct[T any] struct {
	gormDB *gorm.DB
}

func MakeRepository[T any](gormDB *gorm.DB) Repository[T] {
	return repositoryStruct[T]{gormDB: gormDB}
}

func (repository repositoryStruct[T]) Execute(query Query[T]) (*T, error) {
	result, err := query(repository.gormDB)
	return result, err
}
