package reads

import (
	"errors"

	"context"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func StreamProcessingReportById(id uint) repository.QueryOne[entities.StreamProcessingReport] {
	return func(ctx context.Context, gormDB *gorm.DB) (*entities.StreamProcessingReport, error) {
		var report entities.StreamProcessingReport
		result := gormDB.First(&report, id)

		if result.Error == nil {
			return &report, nil
		}

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
}
