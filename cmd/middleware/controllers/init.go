package controllers

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers/notifications"
	"github.com/labstack/echo/v4"
)

func Bind(path string, e *echo.Echo) *echo.Group {
	group := e.Group(path)
	notifications.RegisterHandlers(group)

	return group
}
