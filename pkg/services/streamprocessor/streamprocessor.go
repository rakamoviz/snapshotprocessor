package streamprocessor

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities/streamprocessingstatus"
	"bufio"
	"fmt"
	"gorm.io/gorm"
	"log"
	"sync"
)

const CHUNK_LEN = 5

type (
	ProcessLine func(line string, gormDB *gorm.DB) error
	OpenScanner func(path string) (*bufio.Scanner, error)
)

type StreamProcessor interface {
	Run(
		path string,
		reportReference string,
		ignoreFirst bool,
		reportsCh chan<- entities.StreamProcessingReport,
		errorsCh chan<- error,
		processLine ProcessLine,
	)
}

type streamProcessor struct {
	gormDB      *gorm.DB
	openScanner OpenScanner
}

func New(gormDB *gorm.DB, openScanner OpenScanner) *streamProcessor {
	return &streamProcessor{gormDB: gormDB, openScanner: openScanner}
}

func (sproc *streamProcessor) processChunk(
	ignoreFirst bool, chunk []string, chunkId int, processLine ProcessLine,
	report entities.StreamProcessingReport,
) entities.ChunkProcessingReport {
	lineNumberOffset := uint32((chunkId * CHUNK_LEN) + 1)
	if ignoreFirst {
		lineNumberOffset += uint32(1)
	}

	errCount := 0
	for i, line := range chunk {
		lineNumber := uint32(lineNumberOffset) + uint32(i)
		err := processLine(line, sproc.gormDB)
		if err != nil {
			lpError := entities.LineProcessingError{
				LineNumber: lineNumber,
				Line:       line,
				ReportID:   report.ID,
				Error:      err.Error(),
			}

			err := sproc.gormDB.Create(&lpError).Error
			if err != nil {
				log.Printf("Failed to save LineProcessingError %v for %v\n", lpError, err.Error())
			}

			errCount++
		}
	}

	return entities.ChunkProcessingReport{SuccessCount: uint32(len(chunk) - errCount), ErrorsCount: uint32(errCount)}
}

func (sproc *streamProcessor) Run(
	path string,
	reportReference string,
	ignoreFirst bool,
	reportsCh chan<- entities.StreamProcessingReport,
	errorsCh chan<- error,
	processLine ProcessLine,
) {
	fmt.Println("TOLOL 1")
	var procChunksGathererWG sync.WaitGroup

	chunkProcessingReportsCh := make(chan entities.ChunkProcessingReport)

	streamProcessingReport := entities.StreamProcessingReport{
		Reference: reportReference,
	}
	streamProcessingReport.Path = path
	streamProcessingReport.ChunkProcessingReport = entities.ChunkProcessingReport{
		SuccessCount: 0,
		ErrorsCount:  0,
	}

	fmt.Println("TOLOL 2.0.A")
	if reportsCh != nil {
		reportsCh <- streamProcessingReport
	}
	fmt.Println("TOLOL 2.0")

	err := sproc.gormDB.Create(&streamProcessingReport).Error
	if err != nil {
		if errorsCh != nil {
			errorsCh <- err
		}
		return
	}

	go func() {
		for chunkProcessingReport := range chunkProcessingReportsCh {
			streamProcessingReport.SuccessCount += chunkProcessingReport.SuccessCount
			streamProcessingReport.ErrorsCount += chunkProcessingReport.ErrorsCount

			if reportsCh != nil {
				reportsCh <- streamProcessingReport
			}

			err = sproc.gormDB.Save(&streamProcessingReport).Error
			if err != nil {
				log.Printf("Failed updating StreamProcessingReport %v\n", streamProcessingReport) //TODO: provide more detail
				if errorsCh != nil {
					errorsCh <- err
				}
			}

			procChunksGathererWG.Done()
		}
	}()

	fmt.Println("TOLOL 2")
	var chunk []string
	chunkId := 0
	chunkLineOffset := 0

	scanner, err := sproc.openScanner(path)
	if err != nil {
		log.Printf("Failed opening scanner StreamProcessingReport %v\n", streamProcessingReport) //TODO: provide more detail
		if errorsCh != nil {
			errorsCh <- err
		}
		return
	}

	streamProcessingReport.Status = streamprocessingstatus.Running

	reportsCh <- streamProcessingReport

	err = sproc.gormDB.Save(&streamProcessingReport).Error
	if err != nil {
		log.Printf("Failed updating StreamProcessingReport %v\n", streamProcessingReport) //TODO: provide more detail
		if errorsCh != nil {
			errorsCh <- err
		}
		return
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
					chunkProcessingReportsCh <- sproc.processChunk(
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
			chunkProcessingReportsCh <- sproc.processChunk(
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

	if reportsCh != nil {
		reportsCh <- streamProcessingReport
	}

	err = sproc.gormDB.Save(streamProcessingReport).Error
	if err != nil {
		log.Printf("Failed updating StreamProcessingReport %v\n", streamProcessingReport) //TODO: provide more detail
		if errorsCh != nil {
			errorsCh <- err
		}
		return
	}
}
