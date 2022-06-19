package provider1

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"gorm.io/gorm"
)

func ProcessFormat1(line string, gormDB *gorm.DB) error {
	columns := strings.Split(line, ",")
	if len(columns) < 6 {
		return fmt.Errorf("%v: has less than 6 columns", line)
	}

	return gormDB.Transaction(func(tx *gorm.DB) error {
		cluster := entities.Cluster{
			Code: columns[0][1 : len(columns[0])-1],
		}

		err := gormDB.FirstOrCreate(&cluster, cluster).Error
		if err != nil {
			return err
		}

		node := entities.Node{
			Code:      columns[1][1 : len(columns[1])-1],
			ClusterID: cluster.Code,
		}

		err = gormDB.FirstOrCreate(&node, node).Error
		if err != nil {
			return err
		}

		timestamp, err := strconv.ParseInt(strings.Trim(columns[2], " "), 0, 64)
		if err != nil {
			return err
		}

		nodeStatus := entities.NodeStatus{
			NodeID: node.Code,
			Time:   time.Unix(timestamp, 0),
		}

		err = gormDB.Create(&nodeStatus).Error
		if err != nil {
			return err
		}

		return nil
	})
}
