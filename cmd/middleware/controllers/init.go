package controllers

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers/notifications"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/workers/streamprocessing"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Bind(group *echo.Group)
}

func Setup(g *echo.Group, streamProcessingWorker streamprocessing.Worker) {
	notificationsHandler := notifications.New(streamProcessingWorker)
	notificationsGroup := g.Group("/notifications")
	notificationsHandler.Bind(notificationsGroup)
}
