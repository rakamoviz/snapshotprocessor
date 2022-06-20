package lists

import (
	"errors"

	"context"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/misc"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func StreamProcessingReports(queryParams misc.ListQueryParams) repository.Query[entities.StreamProcessingReport] {
	return func(ctx context.Context, gormDB *gorm.DB) ([]entities.StreamProcessingReport, error) {
		var reports []entities.StreamProcessingReport

		query := misc.ListQueryParamsToQuery(
			ctx,
			gormDB,
			queryParams,
		)
		result := query.Find(&reports)

		if result.Error == nil {
			return reports, nil
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
}
