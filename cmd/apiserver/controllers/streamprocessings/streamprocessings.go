package streamprocessings

import (
	"github.com/labstack/echo/v4"
	"github.com/rakamoviz/snapshotprocessor/cmd/apiserver/middlewares"
	"github.com/rakamoviz/snapshotprocessor/pkg/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/repository"
	"github.com/rakamoviz/snapshotprocessor/pkg/scheduler"
	"github.com/rakamoviz/snapshotprocessor/pkg/scheduler/handlers"
)

type controller struct {
	streamProcessingScheduler        scheduler.Client[handlers.StreamProcessingJobData]
	apiKeyCheckMiddleware            *middlewares.ApiKeyCheck
	streamProcessingReportRepository repository.Repository[entities.StreamProcessingReport]
	lineProcessingErrorRepository    repository.Repository[entities.LineProcessingError]
}

type receiveResponse struct {
	JobID           string `json:"job_id"`
	ReportReference string `json:"report_reference"`
}

func New(
	streamProcessingScheduler scheduler.Client[handlers.StreamProcessingJobData],
	apiKeyCheckMiddleware *middlewares.ApiKeyCheck,
	streamProcessingReportRepository repository.Repository[entities.StreamProcessingReport],
	lineProcessingErrorRepository repository.Repository[entities.LineProcessingError],
) *controller {
	return &controller{
		streamProcessingScheduler:        streamProcessingScheduler,
		apiKeyCheckMiddleware:            apiKeyCheckMiddleware,
		streamProcessingReportRepository: streamProcessingReportRepository,
		lineProcessingErrorRepository:    lineProcessingErrorRepository,
	}
}

func (c *controller) Bind(group *echo.Group) {
	group.POST("", func(ctx echo.Context) error { return c.enqueueStreamProcessing(ctx) }, c.apiKeyCheckMiddleware.Process)
	group.GET("/:id", func(ctx echo.Context) error { return c.getByID(ctx) })
	group.GET("", func(ctx echo.Context) error { return c.list(ctx) })
	//group.GET("/:id/errors", func(ctx echo.Context) error { return c.getByID(ctx) })
}
