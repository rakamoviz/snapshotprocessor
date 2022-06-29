package upserts

import (
	"context"

	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Node(node entities.Node) repository.QueryOne[entities.Node] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.Node, error) {
		err := gormDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&node).Error
		if err != nil {
			return nil, err
		}

		return &node, err
	}
}
