package streamprocessings

import (
	"net/http"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities/lists"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/misc"
	"github.com/labstack/echo/v4"
)

func (c *controller) list(ctx echo.Context) error {
	listQueryParams, err := misc.BindListQueryParams(ctx)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Bad Request")
	}

	streamProcessingReports, err := c.streamProcessingReportRepository.Execute(
		ctx.Request().Context(),
		lists.StreamProcessingReports(listQueryParams),
	)

	if err == nil {
		return ctx.JSON(http.StatusOK, streamProcessingReports)
	}

	return err
}
