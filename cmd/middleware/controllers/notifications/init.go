package notifications

import (
	"github.com/labstack/echo/v4"
)

func Bind(g *echo.Group) {
	g.GET("/", receive)
}
