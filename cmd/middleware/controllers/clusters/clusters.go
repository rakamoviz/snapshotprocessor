package clusters

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type controller struct {
	gormDB *gorm.DB
}

type readResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

func New(gormDB *gorm.DB) *controller {
	return &controller{gormDB: gormDB}
}

func (c *controller) Bind(group *echo.Group) {
	group.GET("/:id", func(ctx echo.Context) error { return c.getByID(ctx) })
}
