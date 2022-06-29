package createfirsts

import (
	"context"

	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func Node(node entities.Node) repository.QueryOne[entities.Node] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.Node, error) {
		err := gormDB.Where(entities.Node{Code: node.Code}).FirstOrCreate(&node).Error
		if err != nil {
			return nil, err
		}

		return &node, nil
	}
}
