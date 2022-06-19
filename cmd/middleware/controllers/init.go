package controllers

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers/clusters"
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers/streamprocessings"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/scheduler"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler interface {
	Bind(group *echo.Group)
}

func Setup(
	g *echo.Group, gormDB *gorm.DB,
	streamProcessingScheduler scheduler.Client[handlers.StreamProcessingJobData],
) {
	streamprocessingsHandler := streamprocessings.New(gormDB, streamProcessingScheduler)
	streamprocessingsGroup := g.Group("/streamprocessings")
	streamprocessingsHandler.Bind(streamprocessingsGroup)

	clustersHandler := clusters.New(gormDB)
	clustersGroup := g.Group("/clusters")
	clustersHandler.Bind(clustersGroup)
}
