package models

import (
	"time"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models/streamprocessingstatus"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Cluster struct {
	gorm.Model
	Code string          `gorm:"column:code;uniqueIndex;"`
	Lat  decimal.Decimal `gorm:"not null;type=decimal(3,15);"`
	Lng  decimal.Decimal `gorm:"not null;type=decimal(3,15);"`
}

type Node struct {
	gorm.Model
	Code      string `gorm:"column:code;uniqueIndex;"`
	ClusterID string
	Cluster   Cluster `gorm:"references:Code;not null;"`
}

type NodeStatus struct {
	gorm.Model
	NodeID      string
	Node        Node `gorm:"references:Code;not null;"`
	CpuUsage    decimal.Decimal
	DiskUsage   uint64
	MemoryUsage uint64
	Time        time.Time
}

type ChunkProcessingReport struct {
	SuccessCount uint32
	ErrorsCount  uint32
}

type StreamProcessingReport struct {
	gorm.Model
	ChunkProcessingReport `gorm:"embedded"`
	Path                  string
	Status                streamprocessingstatus.Enum
}

type LineProcessingError struct {
	gorm.Model
	LineNumber uint32
	Line       string
	ReportID   uint
	Report     StreamProcessingReport
	Error      string
}
