package entities

import (
	"time"

	pkgentities "bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"github.com/shopspring/decimal"
)

type Cluster struct {
	pkgentities.Entity
	Code string `gorm:"column:code;uniqueIndex;"`
}

type Node struct {
	pkgentities.Entity
	Code      string `gorm:"column:code;uniqueIndex;"`
	ClusterID string
	Cluster   Cluster `gorm:"references:Code;not null;"`
}

type NodeStatus struct {
	pkgentities.Entity
	NodeID      string
	Node        Node `gorm:"references:Code;not null;"`
	CpuUsage    decimal.Decimal
	DiskUsage   uint64
	MemoryUsage uint64
	Time        time.Time
}
