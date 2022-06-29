package lists

import (
	"errors"

	"context"

	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/misc"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func Nodes(queryParams misc.ListQueryParams) repository.Query[entities.Node] {
	return func(ctx context.Context, gormDB *gorm.DB) ([]entities.Node, error) {
		var nodes []entities.Node

		query := misc.ListQueryParamsToQuery(
			ctx,
			gormDB,
			queryParams,
		)
		result := query.Find(&nodes)

		if result.Error == nil {
			return nodes, nil
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
}
