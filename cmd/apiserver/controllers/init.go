package controllers

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/apiserver/controllers/clusters"
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/apiserver/controllers/streamprocessings"
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/apiserver/middlewares"
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
	apiKeyCheckMiddleware *middlewares.ApiKeyCheck,
) {
	streamprocessingsHandler := streamprocessings.New(gormDB, streamProcessingScheduler, apiKeyCheckMiddleware)
	streamprocessingsGroup := g.Group("/streamprocessings")
	streamprocessingsHandler.Bind(streamprocessingsGroup)

	clustersHandler := clusters.New(gormDB)
	clustersGroup := g.Group("/clusters")
	clustersHandler.Bind(clustersGroup)
}
