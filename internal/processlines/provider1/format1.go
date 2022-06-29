package provider1

import (
	"fmt"
	"strings"

	"context"

	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/internal/entities/createfirsts"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"gorm.io/gorm"
)

func ProcessFormat1(
	ctx context.Context, line string, gormDB *gorm.DB,
) error {
	columns := strings.Split(line, ",")
	if len(columns) < 6 {
		return fmt.Errorf("%v: has less than 6 columns", line)
	}

	clusterRepository := repository.New[entities.Cluster](gormDB)
	nodeRepository := repository.New[entities.Node](gormDB)

	return gormDB.Transaction(func(tx *gorm.DB) error {
		cluster := entities.Cluster{
			Code: strings.Trim(columns[0], "\""),
		}

		_, err := clusterRepository.ExecuteOne(ctx, createfirsts.Cluster(cluster))
		if err != nil {
			return err
		}

		fmt.Println("Cluster is ", cluster)

		node := entities.Node{
			Code:    strings.Trim(columns[1], "\""),
			Cluster: cluster,
		}

		_, err = nodeRepository.ExecuteOne(ctx, createfirsts.Node(node))
		if err != nil {
			return err
		}

		return nil
	})
}
