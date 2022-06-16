package creates

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreateCluster(code string) repository.Query[models.Cluster] {
	return func(gormDB *gorm.DB) (*models.Cluster, error) {
		var err error

		pCluster := &models.Cluster{Code: code}

		err = gormDB.Clauses(clause.OnConflict{UpdateAll: true}).Create(pCluster).Error
		return pCluster, err
	}
}
