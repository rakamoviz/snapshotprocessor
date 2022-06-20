package controllers

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/apiserver/controllers/clusters"
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/apiserver/controllers/streamprocessings"
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/apiserver/middlewares"
	internalentities "bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	pkgentities "bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
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
	clusterRepository repository.Repository[internalentities.Cluster],
	nodeRepository repository.Repository[internalentities.Node],
	nodeStatusRepository repository.Repository[internalentities.NodeStatus],
	streamProcessingReportRepository repository.Repository[pkgentities.StreamProcessingReport],
	lineProcessingErrorRepository repository.Repository[pkgentities.LineProcessingError],
) {
	streamprocessingsHandler := streamprocessings.New(
		streamProcessingScheduler, apiKeyCheckMiddleware,
		streamProcessingReportRepository, lineProcessingErrorRepository,
	)
	streamprocessingsGroup := g.Group("/streamprocessings")
	streamprocessingsHandler.Bind(streamprocessingsGroup)

	clustersHandler := clusters.New(
		clusterRepository,
	)
	clustersGroup := g.Group("/clusters")
	clustersHandler.Bind(clustersGroup)
}
