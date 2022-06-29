package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/rakamoviz/snapshotprocessor/cmd/apiserver/controllers/clusters"
	"github.com/rakamoviz/snapshotprocessor/cmd/apiserver/controllers/nodes"
	"github.com/rakamoviz/snapshotprocessor/cmd/apiserver/controllers/nodestatuses"
	"github.com/rakamoviz/snapshotprocessor/cmd/apiserver/controllers/streamprocessings"
	"github.com/rakamoviz/snapshotprocessor/cmd/apiserver/middlewares"
	internalentities "github.com/rakamoviz/snapshotprocessor/internal/entities"
	pkgentities "github.com/rakamoviz/snapshotprocessor/pkg/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"github.com/rakamoviz/snapshotprocessor/pkg/scheduler"
	"github.com/rakamoviz/snapshotprocessor/pkg/scheduler/handlers"
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

	nodesHandler := nodes.New(
		nodeRepository,
	)
	nodesGroup := g.Group("/nodes")
	nodesHandler.Bind(nodesGroup)

	nodeStatusesHandler := nodestatuses.New(
		nodeStatusRepository,
	)
	nodeStatusesGroup := g.Group("/nodestatuses")
	nodeStatusesHandler.Bind(nodeStatusesGroup)
}
