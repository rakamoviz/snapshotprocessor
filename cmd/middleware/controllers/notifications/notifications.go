package notifications

import (
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities/streamprocessingstatus"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/workers/streamprocessing"
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

	streamProcessingReportCh := c.streamProcessingWorker.AppendJob(path, true, processLine)
	streamProcessingReport := <-streamProcessingReportCh

	ctx.JSON(http.StatusOK, &postResponse{
		ReportID: streamProcessingReport.ID,
	})

	for {
		report := <-streamProcessingReportCh
		fmt.Printf("Report. Status: %v, Success:%v, Errors:%v\n", report.Status, report.SuccessCount, report.ErrorsCount)
		if report.Status != streamprocessingstatus.Running && report.Status != streamprocessingstatus.Undefined {
			break
		}
	}
	return nil
}
