package createfirsts

import (
	"context"

	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func Cluster(cluster entities.Cluster) repository.QueryOne[entities.Cluster] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.Cluster, error) {
		err := gormDB.Where(entities.Cluster{Code: cluster.Code}).FirstOrCreate(&cluster).Error
		if err != nil {
			return nil, err
		}

		return &cluster, nil
	}
}
