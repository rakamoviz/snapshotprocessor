package notifications

import (
	"github.com/labstack/echo/v4"
)

func registerHandlers(g *echo.Group) {
	g.POST("/", receive)
}
