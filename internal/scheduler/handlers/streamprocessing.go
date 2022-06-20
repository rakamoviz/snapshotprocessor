package handlers

import (
	"context"
	"fmt"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	//"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities/streamprocessingstatus"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/services/streamprocessor"
	"gorm.io/gorm"
)

type StreamProcessingJobData struct {
	Provider        string `json:"provider,omitempty"`
	Format          string `json:"format,omitempty"`
	Path            string `json:"path,omitempty"`
	ReportReference string `json:"report_reference,omitempty"`
}

type streamProcessing[T StreamProcessingJobData] struct {
	gormDB          *gorm.DB
	streamProcessor streamprocessor.StreamProcessor
	processLines    map[string]map[string]streamprocessor.ProcessLine
}

func NewStreamProcessing(
	gormDB *gorm.DB,
	streamProcessor streamprocessor.StreamProcessor,
	processLines map[string]map[string]streamprocessor.ProcessLine,
) *streamProcessing[StreamProcessingJobData] {
	h := streamProcessing[StreamProcessingJobData]{
		gormDB:          gormDB,
		streamProcessor: streamProcessor,
		processLines:    processLines,
	}

	return &h
}

func (h *streamProcessing[T]) Handle(ctx context.Context, jobData StreamProcessingJobData) error {
	streamProcessingReport := entities.StreamProcessingReport{
		Reference: jobData.ReportReference,
	}
	streamProcessingReport.Path = jobData.Path
	streamProcessingReport.Provider = jobData.Provider
	streamProcessingReport.Format = jobData.Format
	streamProcessingReport.ChunkProcessingReport = entities.ChunkProcessingReport{
		SuccessCount: 0,
		ErrorsCount:  0,
	}

	formats, ok := h.processLines[jobData.Provider]
	if !ok {
		err := fmt.Errorf("No formats registered for provider %s", jobData.Provider)
		streamProcessingReport.Error = err.Error()
		errSavingStreamProcessingReport := h.gormDB.Create(&streamProcessingReport)
		if errSavingStreamProcessingReport != nil {
			fmt.Println(err.Error())
		}
		return err
	}

	processLine, ok := formats[jobData.Format]
	if !ok {
		err := fmt.Errorf("No formats registered for provider %s, format %s", jobData.Provider, jobData.Format)
		streamProcessingReport.Error = err.Error()
		errSavingStreamProcessingReport := h.gormDB.Create(&streamProcessingReport)
		if errSavingStreamProcessingReport != nil {
			fmt.Println(err.Error())
		}
		return err
	}

	streamProcessingReportCh := make(chan entities.StreamProcessingReport)
	errorsCh := make(chan error)

	/*
		go func() {
			for {
				report := <-streamProcessingReportCh
				fmt.Printf("Report. Status: %v, Success:%v, Errors:%v\n", report.Status, report.SuccessCount, report.ErrorsCount)
				if report.Status != streamprocessingstatus.Running && report.Status != streamprocessingstatus.Undefined {
					break
				}
			}
		}()
	*/

	h.streamProcessor.Run(
		ctx,
		jobData.Path,
		jobData.ReportReference,
		true,
		streamProcessingReportCh,
		errorsCh,
		processLine,
	)

	return nil
}
