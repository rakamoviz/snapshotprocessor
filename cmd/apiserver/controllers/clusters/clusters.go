package clusters

import (
	"github.com/labstack/echo/v4"
	"github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
)

type controller struct {
	clusterRepository repository.Repository[entities.Cluster]
}

func New(clusterRepository repository.Repository[entities.Cluster]) *controller {
	return &controller{clusterRepository: clusterRepository}
}

func (c *controller) Bind(group *echo.Group) {
	group.GET("/:id", func(ctx echo.Context) error { return c.getByID(ctx) })
	group.GET("", func(ctx echo.Context) error { return c.list(ctx) })
}
