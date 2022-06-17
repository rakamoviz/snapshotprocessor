package streamprocessor

import (
	"bufio"
	"log"
	"sync"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models/streamprocessingstatus"
	"gorm.io/gorm"
)

const CHUNK_LEN = 5

type (
	ProcessLine func(line string, gormDB *gorm.DB) error
	OpenScanner func(path string) (*bufio.Scanner, error)
)

type StreamProcessor interface {
	Run(
		path string,
		ignoreFirst bool,
		reportCh chan models.StreamProcessingReport,
		processLine ProcessLine,
	) error
}

type streamProcessorStruct struct {
	gormDB      *gorm.DB
	openScanner OpenScanner
}

func MakeStreamProcessor(gormDB *gorm.DB, openScanner OpenScanner) StreamProcessor {
	return streamProcessorStruct{gormDB: gormDB, openScanner: openScanner}
}

func (streamProcessor streamProcessorStruct) processChunk(
	ignoreFirst bool, chunk []string, chunkId int, processLine ProcessLine,
	report models.StreamProcessingReport,
) models.ChunkProcessingReport {
	lineNumberOffset := uint32((chunkId * CHUNK_LEN) + 1)
	if ignoreFirst {
		lineNumberOffset += uint32(1)
	}

	errCount := 0
	for i, line := range chunk {
		lineNumber := uint32(lineNumberOffset) + uint32(i)
		err := processLine(line, streamProcessor.gormDB)
		if err != nil {
			lpError := models.LineProcessingError{
				LineNumber: lineNumber,
				Line:       line,
				ReportID:   report.ID,
				Error:      err.Error(),
			}

			err := streamProcessor.gormDB.Create(&lpError).Error
			if err != nil {
				log.Printf("Failed to save LineProcessingError %v for %v\n", lpError, err.Error())
			}

			errCount++
		}
	}

	return models.ChunkProcessingReport{SuccessCount: uint32(len(chunk) - errCount), ErrorsCount: uint32(errCount)}
}

func (streamProcessor streamProcessorStruct) Run(
	path string,
	ignoreFirst bool,
	reportCh chan models.StreamProcessingReport,
	processLine ProcessLine,
) error {
	var procChunksGathererWG sync.WaitGroup

	chunkProcessingReportsCh := make(chan models.ChunkProcessingReport)

	streamProcessingReport := models.StreamProcessingReport{}
	streamProcessingReport.Path = path
	streamProcessingReport.ChunkProcessingReport = models.ChunkProcessingReport{
		SuccessCount: 0,
		ErrorsCount:  0,
	}

	reportCh <- streamProcessingReport

	err := streamProcessor.gormDB.Create(&streamProcessingReport).Error
	if err != nil {
		return err
	}

	go func() {
		for chunkProcessingReport := range chunkProcessingReportsCh {
			streamProcessingReport.SuccessCount += chunkProcessingReport.SuccessCount
			streamProcessingReport.ErrorsCount += chunkProcessingReport.ErrorsCount

			reportCh <- streamProcessingReport

			err = streamProcessor.gormDB.Save(&streamProcessingReport).Error
			if err != nil {
				log.Printf("Failed updating StreamProcessingReport %v\n", streamProcessingReport) //TODO: provide more detail
			}

			procChunksGathererWG.Done()
		}
	}()

	var chunk []string
	chunkId := 0
	chunkLineOffset := 0

	scanner, err := streamProcessor.openScanner(path)
	if err != nil {
		log.Printf("Failed opening scanner StreamProcessingReport %v\n", streamProcessingReport) //TODO: provide more detail
		return err
	}

	streamProcessingReport.Status = streamprocessingstatus.Running

	reportCh <- streamProcessingReport

	err = streamProcessor.gormDB.Save(&streamProcessingReport).Error
	if err != nil {
		log.Printf("Failed updating StreamProcessingReport %v\n", streamProcessingReport) //TODO: provide more detail
		return err
	}

	for scanner.Scan() {
		if ignoreFirst {
			ignoreFirst = false
			continue
		}

		line := scanner.Text()

		if chunkLineOffset%CHUNK_LEN == 0 {
			if chunkId > 0 {
				procChunksGathererWG.Add(1)
				go func(chunk []string, chunkId int) {
					chunkProcessingReportsCh <- streamProcessor.processChunk(
						ignoreFirst, chunk, chunkId, processLine,
						streamProcessingReport,
					)
				}(chunk, chunkId)
			}

			chunkLineOffset = 0
			chunk = make([]string, 0, CHUNK_LEN)
			chunkId += 1
		}

		chunk = append(chunk, line)
		chunkLineOffset += 1
	}

	if chunkId > 0 {
		procChunksGathererWG.Add(1)
		go func(chunk []string, chunkId int) {
			chunkProcessingReportsCh <- streamProcessor.processChunk(
				ignoreFirst,
				chunk, chunkId, processLine,
				streamProcessingReport,
			)
		}(chunk, chunkId)
	}

	procChunksGathererWG.Wait()

	if streamProcessingReport.ErrorsCount == 0 {
		streamProcessingReport.Status = streamprocessingstatus.Success
	} else if streamProcessingReport.SuccessCount == 0 {
		streamProcessingReport.Status = streamprocessingstatus.Failed
	} else {
		streamProcessingReport.Status = streamprocessingstatus.Partial
	}

	reportCh <- streamProcessingReport

	err = streamProcessor.gormDB.Save(streamProcessingReport).Error
	if err != nil {
		log.Printf("Failed updating StreamProcessingReport %v\n", streamProcessingReport) //TODO: provide more detail
	}

	return err

}
