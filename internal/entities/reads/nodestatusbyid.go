package reads

import (
	"errors"

	"context"

	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func NodeStatusByID(id uint) repository.QueryOne[entities.NodeStatus] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.NodeStatus, error) {
		var nodeStatus entities.NodeStatus
		result := gormDB.First(&nodeStatus, id)

		if result.Error == nil {
			return &nodeStatus, nil
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
}
