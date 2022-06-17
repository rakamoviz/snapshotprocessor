package creates

import (
	"fmt"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/repository"
	"gorm.io/gorm"
)

func CreateNode(
	code string, clusterID string,
) repository.Query[models.Node] {
	return func(gormDB *gorm.DB) (*models.Node, error) {
		var err error

		pNode := &models.Node{
			Code:      code,
			ClusterID: clusterID,
		}

		fmt.Println("-----------------------------------------------------------")
		fmt.Println(pNode)
		fmt.Println("-----------------------------------------------------------")
		err = gormDB.Create(pNode).Error

		return pNode, err
	}
}
