package controllers

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers/notifications"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/streamprocessor"
	"github.com/labstack/echo/v4"
)

type Handler interface {
	Bind(group *echo.Group)
}

func Setup(g *echo.Group, streamProcessor streamprocessor.StreamProcessor) {
	notificationsHandler := notifications.NewHandler(streamProcessor)
	notificationsGroup := g.Group("/notifications")
	notificationsHandler.Bind(notificationsGroup)
}
