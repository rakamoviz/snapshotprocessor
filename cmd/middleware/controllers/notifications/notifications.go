package notifications

import (
	"net/http"
	"strings"

	"fmt"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/scheduler"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type controller struct {
	streamProcessingScheduler scheduler.Client[handlers.StreamProcessingJobData]
}

type receiveResponse struct {
	JobID           string `json:"job_id"`
	ReportReference string `json:"report_reference"`
}

func New(streamProcessingScheduler scheduler.Client[handlers.StreamProcessingJobData]) *controller {
	return &controller{streamProcessingScheduler: streamProcessingScheduler}
}

func (c *controller) Bind(group *echo.Group) {
	group.POST("", func(ctx echo.Context) error { return c.receiveNotification(ctx) })
}

func (c *controller) receiveNotification(ctx echo.Context) error {
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

	fmt.Println(jobID)

	if err != nil {
		ctx.String(http.StatusInternalServerError, "can't register job")
	}

	ctx.JSON(http.StatusOK, &receiveResponse{
		JobID:           jobID,
		ReportReference: reportReference.String(),
	})

	return nil
}
