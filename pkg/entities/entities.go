package entities

import (
	"github.com/rakamoviz/snapshotprocessor/pkg/entities/streamprocessingstatus"
	"gorm.io/gorm"
)

type Entity gorm.Model

type ChunkProcessingReport struct {
	SuccessCount uint32 `gorm:"not null"`
	ErrorsCount  uint32 `gorm:"not null"`
}

type StreamProcessingReport struct {
	Entity
	Reference             string `gorm:"uniqueIndex;not null"`
	Provider              string
	Format                string
	Path                  string                      `gorm:"not null"`
	Status                streamprocessingstatus.Enum `gorm:"not null"`
	ChunkProcessingReport `gorm:"embedded"`
	Error                 string
}

type LineProcessingError struct {
	Entity
	LineNumber uint32
	Line       string `gorm:"not null"`
	Error      string
	ReportID   uint                   `gorm:"not null"`
	Report     StreamProcessingReport `gorm:"not null"`
}
