package nodes

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rakamoviz/snapshotprocessor/internal/entities/reads"
	"gorm.io/gorm"
)

func (c *controller) getByID(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Incorrect format of id query parameter")
	}

	node, err := c.nodeRepository.ExecuteOne(
		ctx.Request().Context(),
		reads.NodeByID(uint(id)),
	)

	if err == nil {
		return ctx.JSON(http.StatusOK, node)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx.String(http.StatusNotFound, fmt.Sprintf("Node with id %d not found", id))
	}

	return err
}
