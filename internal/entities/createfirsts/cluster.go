package createfirsts

import (
	"context"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func Cluster_First(code string) repository.QueryOne[entities.Cluster] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.Cluster, error) {
		cluster := entities.Cluster{Code: code}

		err := gormDB.FirstOrCreate(&cluster, cluster).Error
		if err != nil {
			return nil, err
		}

		return &cluster, nil
	}
}
