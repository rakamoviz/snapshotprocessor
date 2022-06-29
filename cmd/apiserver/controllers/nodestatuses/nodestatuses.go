package nodestatuses

import (
	"github.com/labstack/echo/v4"
	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
)

type controller struct {
	nodeStatusRepository repository.Repository[entities.NodeStatus]
}

func New(nodeStatusRepository repository.Repository[entities.NodeStatus]) *controller {
	return &controller{nodeStatusRepository: nodeStatusRepository}
}

func (c *controller) Bind(group *echo.Group) {
	group.GET("/:id", func(ctx echo.Context) error { return c.getByID(ctx) })
}
