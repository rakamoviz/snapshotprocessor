package notifications

import (
	"net/http"
	"strings"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/workers/streamprocessing"
	"fmt"
	"github.com/labstack/echo/v4"
)

type controller struct {
	streamProcessingWorker streamprocessing.Worker
}

type postResponse struct {
	ReportID uint `json:"report_id"`
}

func New(streamProcessingWorker streamprocessing.Worker) *controller {
	return &controller{streamProcessingWorker: streamProcessingWorker}
}

func (c *controller) Bind(group *echo.Group) {
	group.POST("", func(ctx echo.Context) error { return c.post(ctx) })
}

func (c *controller) post(ctx echo.Context) error {
	path := ctx.QueryParam("path")
	if strings.Trim(path, " ") == "" {
		ctx.String(http.StatusBadRequest, "missing path query parameter")
		return nil
	}

	fmt.Println("1")
	streamProcessingReportCh := c.streamProcessingWorker.AppendJob(path, true, processLine)
	fmt.Println("2")
	streamProcessingReport := <-streamProcessingReportCh
	fmt.Println("3")

	ctx.JSON(http.StatusOK, &postResponse{
		ReportID: streamProcessingReport.ID,
	})

	return nil
}
