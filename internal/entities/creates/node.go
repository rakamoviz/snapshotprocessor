package creates

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func Node(
	code string, clusterID string,
) repository.QueryOne[entities.Node] {
	return func(gormDB *gorm.DB) (*entities.Node, error) {
		node := entities.Node{
			Code:      code,
			ClusterID: clusterID,
		}

		err := gormDB.Create(&node).Error
		if err != nil {
			return nil, err
		}

		return &node, err
	}
}
