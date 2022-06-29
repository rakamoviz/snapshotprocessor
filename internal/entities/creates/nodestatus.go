package creates

import (
	"context"

	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func NodeStatus(nodeStatus entities.NodeStatus) repository.QueryOne[entities.NodeStatus] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.NodeStatus, error) {
		err := gormDB.Create(&nodeStatus).Error
		if err != nil {
			return nil, err
		}

		return &nodeStatus, err
	}
}
