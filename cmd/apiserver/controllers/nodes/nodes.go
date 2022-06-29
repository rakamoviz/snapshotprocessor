package nodes

import (
	"github.com/labstack/echo/v4"
	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
)

type controller struct {
	nodeRepository repository.Repository[entities.Node]
}

func New(nodeRepository repository.Repository[entities.Node]) *controller {
	return &controller{nodeRepository: nodeRepository}
}

func (c *controller) Bind(group *echo.Group) {
	group.GET("/:id", func(ctx echo.Context) error { return c.getByID(ctx) })
	group.GET("", func(ctx echo.Context) error { return c.list(ctx) })
}
