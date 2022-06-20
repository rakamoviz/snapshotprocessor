package clusters

import (
	"errors"
	"fmt"
	"net/http"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (c *controller) getByID(ctx echo.Context) error {
	id := ctx.Param("id")

	var cluster entities.Cluster
	result := c.gormDB.First(&cluster, id)

	if result.Error == nil {
		return ctx.JSON(http.StatusOK, cluster)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("Cluster with id %s not found", id))
	}

	return result.Error
}
