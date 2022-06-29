package clusters

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rakamoviz/snapshotprocessor/internal/entities/lists"
	"github.com/rakamoviz/snapshotprocessor/pkg/misc"
)

func (c *controller) list(ctx echo.Context) error {
	listQueryParams, err := misc.BindListQueryParams(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Bad Request")
	}

	streamProcessingReports, err := c.clusterRepository.Execute(
		ctx.Request().Context(),
		lists.Clusters(listQueryParams),
	)

	if err == nil {
		return ctx.JSON(http.StatusOK, streamProcessingReports)
	}

	return err
}
