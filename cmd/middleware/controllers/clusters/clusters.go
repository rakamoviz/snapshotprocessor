package clusters

import (
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
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

func (c *controller) getByID(ctx echo.Context) error {
	id := ctx.Param("id")

	var cluster entities.Cluster
	result := c.gormDB.First(&cluster, id)

	if result.Error == nil {
		ctx.JSON(http.StatusOK, cluster)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ctx.String(http.StatusNotFound, fmt.Sprintf("Cluster with id %s not found", id))
	}

	return result.Error
}
