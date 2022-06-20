package upserts

import (
	"context"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Cluster(code string) repository.QueryOne[entities.Cluster] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.Cluster, error) {
		cluster := entities.Cluster{Code: code}

		err := gormDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&cluster).Error
		if err != nil {
			return nil, err
		}

		return &cluster, err
	}
}
