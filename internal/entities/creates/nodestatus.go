package creates

import (
	"time"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func CreateNodeStatus(
	nodeID string, time time.Time, cpuUsage decimal.Decimal,
	memoryUsage uint64, diskUsage uint64,
) repository.QueryOne[entities.NodeStatus] {
	return func(gormDB *gorm.DB) (entities.NodeStatus, error) {
		var err error

		nodeStatus := entities.NodeStatus{
			NodeID:      nodeID,
			Time:        time,
			CpuUsage:    cpuUsage,
			MemoryUsage: memoryUsage,
			DiskUsage:   diskUsage,
		}

		err = gormDB.Create(&nodeStatus).Error

		return nodeStatus, err
	}
}
