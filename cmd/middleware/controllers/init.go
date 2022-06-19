package controllers

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers/notifications"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/scheduler"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Bind(group *echo.Group)
}

func Setup(g *echo.Group, streamProcessingScheduler scheduler.Client[handlers.StreamProcessingJobData]) {
	notificationsHandler := notifications.New(streamProcessingScheduler)
	notificationsGroup := g.Group("/notifications")
	notificationsHandler.Bind(notificationsGroup)
}
