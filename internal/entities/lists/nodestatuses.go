package lists

import (
	"errors"

	"context"

	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/misc"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func NodeStatuses(queryParams misc.ListQueryParams) repository.Query[entities.NodeStatus] {
	return func(ctx context.Context, gormDB *gorm.DB) ([]entities.NodeStatus, error) {
		var nodeStatuses []entities.NodeStatus

		query := misc.ListQueryParamsToQuery(
			ctx,
			gormDB,
			queryParams,
		)
		result := query.Find(&nodeStatuses)

		if result.Error == nil {
			return nodeStatuses, nil
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
}
