package clusters

import (
	"net/http"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities/lists"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/misc"
	"github.com/labstack/echo/v4"
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
