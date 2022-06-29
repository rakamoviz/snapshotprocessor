package streamprocessings

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rakamoviz/snapshotprocessor/pkg/scheduler/handlers"
	"github.com/rakamoviz/snapshotprocessor/pkg/services/auth"
)

type enqueueStreamProcessingRequest struct {
	Path   string `json:"path"`
	Format string `json:"format"`
}

func (c *controller) enqueueStreamProcessing(ctx echo.Context) error {
	apiClient := ctx.Get("ApiClient").(auth.ApiClient)

	path := ctx.QueryParam("path")
	if strings.Trim(path, " ") == "" {
		return ctx.String(http.StatusBadRequest, "missing path query parameter")
	}
	format := ctx.QueryParam("format")
	if strings.Trim(format, " ") == "" {
		return ctx.String(http.StatusBadRequest, "missing format query parameter")
	}

	reportReference := uuid.New()
	streamProcessingJobData := handlers.StreamProcessingJobData{
		Provider:        apiClient.Name,
		Format:          format,
		Path:            path,
		ReportReference: reportReference.String(),
	}

	jobID, err := c.streamProcessingScheduler.Schedule(ctx.Request().Context(), streamProcessingJobData, 0)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "Can't register job")
	}

	return ctx.JSON(http.StatusOK, &receiveResponse{
		JobID:           jobID,
		ReportReference: reportReference.String(),
	})
}
