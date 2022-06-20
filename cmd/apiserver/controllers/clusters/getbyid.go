package clusters

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities/reads"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func (c *controller) getByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Incorrect format of id query parameter")
	}

	cluster, err := c.clusterRepository.ExecuteOne(
		ctx.Request().Context(),
		reads.ClusterByID(uint(id)),
	)

	if err == nil {
		return ctx.JSON(http.StatusOK, cluster)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("Cluster with id %d not found", id))
	}

	return err
}
