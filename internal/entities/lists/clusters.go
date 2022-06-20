package lists

import (
	"errors"

	"context"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/misc"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func Clusters(queryParams misc.ListQueryParams) repository.Query[entities.Cluster] {
	return func(ctx context.Context, gormDB *gorm.DB) ([]entities.Cluster, error) {
		var clusters []entities.Cluster

		query := misc.ListQueryParamsToQuery(
			ctx,
			gormDB,
			queryParams,
		)
		result := query.Find(&clusters)

		if result.Error == nil {
			return clusters, nil
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
}
