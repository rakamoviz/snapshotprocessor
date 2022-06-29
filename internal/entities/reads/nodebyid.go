package reads

import (
	"errors"

	"context"

	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func NodeByID(id uint) repository.QueryOne[entities.Node] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.Node, error) {
		var node entities.Node
		result := gormDB.First(&node, id)

		if result.Error == nil {
			return &node, nil
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
}
