package entities

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities/streamprocessingstatus"
	"gorm.io/gorm"
)

type Entity gorm.Model

func (e Entity) HasID() bool {
	return e.ID > 0
}

type ChunkProcessingReport struct {
	SuccessCount uint32 `gorm:"not null"`
	ErrorsCount  uint32 `gorm:"not null"`
}

type StreamProcessingReport struct {
	Entity
	Reference             string `gorm:"uniqueIndex;not null"`
	ChunkProcessingReport `gorm:"embedded"`
	Path                  string                      `gorm:"not null"`
	Status                streamprocessingstatus.Enum `gorm:"not null"`
}

type LineProcessingError struct {
	Entity
	LineNumber uint32
	Line       string `gorm:"not null"`
	Error      string
	ReportID   uint                   `gorm:"not null"`
	Report     StreamProcessingReport `gorm:"not null"`
}
