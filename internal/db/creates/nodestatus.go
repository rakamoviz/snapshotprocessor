package creates

import (
	"time"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func CreateNodeStatus(
	nodeID string, time time.Time, cpuUsage decimal.Decimal,
	memoryUsage uint64, diskUsage uint64,
) repository.Query[models.NodeStatus] {
	return func(gormDB *gorm.DB) (*models.NodeStatus, error) {
		var err error

		pNodeStatus := &models.NodeStatus{
			NodeID:      nodeID,
			Time:        time,
			CpuUsage:    cpuUsage,
			MemoryUsage: memoryUsage,
			DiskUsage:   diskUsage,
		}

		err = gormDB.Create(pNodeStatus).Error

		return pNodeStatus, err
	}
}
