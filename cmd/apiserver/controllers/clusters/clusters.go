package clusters

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"github.com/labstack/echo/v4"
)

type controller struct {
	clusterRepository repository.Repository[entities.Cluster]
}

type readResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
}

func New(clusterRepository repository.Repository[entities.Cluster]) *controller {
	return &controller{clusterRepository: clusterRepository}
}

func (c *controller) Bind(group *echo.Group) {
	group.GET("/:id", func(ctx echo.Context) error { return c.getByID(ctx) })
}
