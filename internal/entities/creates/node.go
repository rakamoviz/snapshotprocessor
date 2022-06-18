package creates

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func CreateNode(
	code string, clusterID string,
) repository.QueryOne[entities.Node] {
	return func(gormDB *gorm.DB) (entities.Node, error) {
		var err error

		node := entities.Node{
			Code:      code,
			ClusterID: clusterID,
		}

		err = gormDB.Create(&node).Error

		return node, err
	}
}
