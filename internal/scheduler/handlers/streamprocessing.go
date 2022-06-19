package handlers

import (
	"fmt"

	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities/streamprocessingstatus"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/services/streamprocessor"
)

type StreamProcessingJobData struct {
	Provider        string `json:"provider,omitempty"`
	Format          string `json:"format,omitempty"`
	Path            string `json:"path,omitempty"`
	ReportReference string `json:"report_reference,omitempty"`
}

type streamProcessing[T StreamProcessingJobData] struct {
	streamProcessor streamprocessor.StreamProcessor
	processLines    map[string]map[string]streamprocessor.ProcessLine
}

func NewStreamProcessing(
	streamProcessor streamprocessor.StreamProcessor,
	processLines map[string]map[string]streamprocessor.ProcessLine,
) *streamProcessing[StreamProcessingJobData] {
	h := streamProcessing[StreamProcessingJobData]{
		streamProcessor: streamProcessor,
		processLines:    processLines,
	}

	return &h
}

func (h *streamProcessing[T]) Handle(jobData StreamProcessingJobData) error {
	fmt.Println("MONYEEETETETETET")
	formats, ok := h.processLines[jobData.Provider]
	if !ok {
		return fmt.Errorf("No formats registered for provider %s", jobData.Provider)
	}

	processLine, ok := formats[jobData.Format]
	if !ok {
		return fmt.Errorf("No formats registered for provider %s", jobData.Provider)
	}

	streamProcessingReportCh := make(chan entities.StreamProcessingReport)
	errorsCh := make(chan error)

	go func() {
		for {
			report := <-streamProcessingReportCh
			fmt.Printf("Report. Status: %v, Success:%v, Errors:%v\n", report.Status, report.SuccessCount, report.ErrorsCount)
			if report.Status != streamprocessingstatus.Running && report.Status != streamprocessingstatus.Undefined {
				break
			}
		}
	}()

	fmt.Println("BANGSAT ")
	h.streamProcessor.Run(
		jobData.Path,
		jobData.ReportReference,
		true,
		streamProcessingReportCh,
		errorsCh,
		processLine,
	)

	fmt.Println("ANJINGSSSSSS ")

	return nil
}
