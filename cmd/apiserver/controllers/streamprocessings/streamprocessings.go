package streamprocessings

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/apiserver/middlewares"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/scheduler"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type controller struct {
	gormDB                    *gorm.DB
	streamProcessingScheduler scheduler.Client[handlers.StreamProcessingJobData]
	apiKeyCheckMiddleware     *middlewares.ApiKeyCheck
}

type receiveResponse struct {
	JobID           string `json:"job_id"`
	ReportReference string `json:"report_reference"`
}

func New(
	gormDB *gorm.DB, streamProcessingScheduler scheduler.Client[handlers.StreamProcessingJobData],
	apiKeyCheckMiddleware *middlewares.ApiKeyCheck,
) *controller {
	return &controller{gormDB: gormDB, streamProcessingScheduler: streamProcessingScheduler, apiKeyCheckMiddleware: apiKeyCheckMiddleware}
}

func (c *controller) Bind(group *echo.Group) {
	group.POST("", func(ctx echo.Context) error { return c.enqueueStreamProcessing(ctx) }, c.apiKeyCheckMiddleware.Process)
	group.GET("/:id", func(ctx echo.Context) error { return c.getByID(ctx) })
	group.GET("", func(ctx echo.Context) error { return c.list(ctx) })
	//group.GET("/:id/errors", func(ctx echo.Context) error { return c.getByID(ctx) })
}
