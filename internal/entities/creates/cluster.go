package creates

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateCluster(code string) repository.QueryOne[entities.Cluster] {
	return func(gormDB *gorm.DB) (entities.Cluster, error) {
		var err error

		cluster := entities.Cluster{Code: code}

		err = gormDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&cluster).Error
		return cluster, err
	}
}
