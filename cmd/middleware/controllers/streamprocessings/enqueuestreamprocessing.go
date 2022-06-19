package streamprocessings

import (
	"net/http"
	"strings"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (c *controller) enqueueStreamProcessing(ctx echo.Context) error {
	path := ctx.QueryParam("path")
	if strings.Trim(path, " ") == "" {
		ctx.String(http.StatusBadRequest, "missing path query parameter")
		return nil
	}
	provider := ctx.QueryParam("provider")
	if strings.Trim(provider, " ") == "" {
		ctx.String(http.StatusBadRequest, "missing provider query parameter")
		return nil
	}
	format := ctx.QueryParam("format")
	if strings.Trim(format, " ") == "" {
		ctx.String(http.StatusBadRequest, "missing format query parameter")
		return nil
	}

	reportReference := uuid.New()
	streamProcessingJobData := handlers.StreamProcessingJobData{
		Provider:        provider,
		Format:          format,
		Path:            path,
		ReportReference: reportReference.String(),
	}
	jobID, err := c.streamProcessingScheduler.Schedule(streamProcessingJobData)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "can't register job")
	}

	ctx.JSON(http.StatusOK, &receiveResponse{
		JobID:           jobID,
		ReportReference: reportReference.String(),
	})

	return nil
}
