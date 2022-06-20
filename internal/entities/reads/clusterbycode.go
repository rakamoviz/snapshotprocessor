package reads

import (
	"errors"

	"context"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func ClusterByCode(code string) repository.QueryOne[entities.Cluster] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.Cluster, error) {
		var cluster entities.Cluster
		result := gormDB.Where(map[string]any{"code": code}).First(&cluster)

		if result.Error == nil {
			return &cluster, nil
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
}
