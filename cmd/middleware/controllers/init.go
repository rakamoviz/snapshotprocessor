package controllers

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers/clusters"
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers/notifications"
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
	notificationsHandler := notifications.New(streamProcessingScheduler)
	notificationsGroup := g.Group("/notifications")
	notificationsHandler.Bind(notificationsGroup)

	clustersHandler := clusters.New(gormDB)
	clustersGroup := g.Group("/clusters")
	clustersHandler.Bind(clustersGroup)
}
