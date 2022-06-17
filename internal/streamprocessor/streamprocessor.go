package streamprocessor

import (
	"bufio"
	"log"
	"sync"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models/streamprocessingstatus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const CHUNK_LEN = 5

type EntitySaveMode uint8

const (
	SaveMode_Noop EntitySaveMode = iota
	SaveMode_InsertIfInexist
	SaveMode_Insert
	SaveMode_Update
	SaveMode_Upsert
)

type (
	ParseLine[T1 any, T2 any, T3 any] func(line string) (*T1, *T2, *T3, error)
	OpenScanner                       func(path string) (*bufio.Scanner, error)
)

type lpErrorsAppender[T any] struct {
}

func (appender lpErrorsAppender[T]) append(
	lpErrors []models.LineProcessingError, failingElts []entityProcessedLineTuple[T],
) []models.LineProcessingError {
	if len(failingElts) == 0 {
		return lpErrors
	}

	newLpErrors := make([]models.LineProcessingError, len(lpErrors)+len(failingElts))
	copy(newLpErrors, lpErrors)
	for idx, failingElt := range failingElts {
		newLpErrors[idx+len(lpErrors)] = models.LineProcessingError{
			ProcessedLine: failingElt.processedLine,
			Error:         failingElt.err.Error(),
		}
	}

	return newLpErrors
}

type StreamProcessor[T1 any, T2 any, T3 any] interface {
	Run(
		path string,
		ignoreFirst bool,
		t1SaveMode EntitySaveMode, t2SaveMode EntitySaveMode, t3SaveMode EntitySaveMode,
		reportCh chan models.StreamProcessingReport,
		parseLine ParseLine[T1, T2, T3],
	) error
}

type streamProcessorStruct[T1 any, T2 any, T3 any] struct {
	gormDB      *gorm.DB
	openScanner OpenScanner
}

type entitiesSaver[T any] struct{ gormDB *gorm.DB }

func (saver entitiesSaver[T]) Perform(saveMode EntitySaveMode, elts []entityProcessedLineTuple[T]) []entityProcessedLineTuple[T] {
	failingElts := make([]entityProcessedLineTuple[T], 0, len(elts))
	switch saveMode {
	case SaveMode_Insert:
		for _, elt := range elts {
			err := saver.gormDB.Create(&elt.entity).Error
			if err != nil {
				failingElt := elt
				failingElt.err = err
				failingElts = append(failingElts, failingElt)
			}
		}
	case SaveMode_InsertIfInexist:
		for _, elt := range elts {

			err := saver.gormDB.FirstOrCreate(&elt.entity, elt.entity).Error
			if err != nil {
				failingElt := elt
				failingElt.err = err
				failingElts = append(failingElts, failingElt)
			}
		}
	case SaveMode_Update:
		for _, elt := range elts {
			err := saver.gormDB.Save(&elt.entity).Error
			if err != nil {
				failingElt := elt
				failingElt.err = err
				failingElts = append(failingElts, failingElt)
			}
		}
	case SaveMode_Upsert:
		for _, elt := range elts {
			err := saver.gormDB.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&elt.entity).Error
			if err != nil {
				failingElt := elt
				failingElt.err = err
				failingElts = append(failingElts, failingElt)
			}
		}
	}
	return failingElts
}

type entityProcessedLineTuple[T any] struct {
	entity        T
	processedLine models.ProcessedLine
	err           error
}

func MakeStreamProcessor[T1 any, T2 any, T3 any](gormDB *gorm.DB, openScanner OpenScanner) StreamProcessor[T1, T2, T3] {
	return streamProcessorStruct[T1, T2, T3]{gormDB: gormDB, openScanner: openScanner}
}

func (streamProcessor streamProcessorStruct[T1, T2, T3]) processChunk(
	ignoreFirst bool, chunk []string, chunkId int, parseLine ParseLine[T1, T2, T3],
	t1SaveMode EntitySaveMode, t2SaveMode EntitySaveMode, t3SaveMode EntitySaveMode,
	report models.StreamProcessingReport,
) models.ChunkProcessingReport {
	lineNumberOffset := uint32((chunkId * CHUNK_LEN) + 1)
	if ignoreFirst {
		lineNumberOffset += uint32(1)
	}

	var elts1 []entityProcessedLineTuple[T1]
	var elts2 []entityProcessedLineTuple[T2]
	var elts3 []entityProcessedLineTuple[T3]

	if t1SaveMode != SaveMode_Noop {
		elts1 = make([]entityProcessedLineTuple[T1], 0, len(chunk))
	}

	if t2SaveMode != SaveMode_Noop {
		elts2 = make([]entityProcessedLineTuple[T2], 0, len(chunk))
	}

	if t3SaveMode != SaveMode_Noop {
		elts3 = make([]entityProcessedLineTuple[T3], 0, len(chunk))
	}

	lineParsingErrors := make([]models.LineProcessingError, 0, len(chunk))

	for i, line := range chunk {
		lineNumber := uint32(lineNumberOffset) + uint32(i)
		processedLine := models.ProcessedLine{
			LineNumber: lineNumber,
			Line:       line,
			Report:     report,
		}

		err := streamProcessor.gormDB.Create(&processedLine).Error

		if err != nil {
			log.Printf("Failed saving line: %v for %v\n", line, err.Error())
			continue
		}

		pEntity1, pEntity2, pEntity3, err := parseLine(line)
		if err != nil {
			lineParsingErrors = append(lineParsingErrors, models.LineProcessingError{
				ProcessedLine: processedLine,
				Error:         err.Error(),
			})
		} else {
			if t1SaveMode != SaveMode_Noop && pEntity1 != nil {
				elts1 = append(elts1, entityProcessedLineTuple[T1]{
					processedLine: processedLine,
					entity:        *pEntity1,
				})
			}
			if t2SaveMode != SaveMode_Noop && pEntity2 != nil {
				elts2 = append(elts2, entityProcessedLineTuple[T2]{
					processedLine: processedLine,
					entity:        *pEntity2,
				})
			}
			if t3SaveMode != SaveMode_Noop && pEntity3 != nil {
				elts3 = append(elts3, entityProcessedLineTuple[T3]{
					processedLine: processedLine,
					entity:        *pEntity3,
				})
			}
		}
	}

	lpErrors := make(
		[]models.LineProcessingError, len(lineParsingErrors),
		len(lineParsingErrors)+len(elts1)+len(elts2)+len(elts3),
	)
	copy(lpErrors, lineParsingErrors)

	savesCount := uint32(0)
	if t1SaveMode != SaveMode_Noop && len(elts1) > 0 {
		saver := entitiesSaver[T1]{gormDB: streamProcessor.gormDB}
		failingElts := saver.Perform(t1SaveMode, elts1)
		lpErrsAppender := lpErrorsAppender[T1]{}
		lpErrors = lpErrsAppender.append(lpErrors, failingElts)
		savesCount += uint32(len(elts1) - len(failingElts))
	}

	if t2SaveMode != SaveMode_Noop && len(elts2) > 0 {
		saver := entitiesSaver[T2]{gormDB: streamProcessor.gormDB}
		failingElts := saver.Perform(t2SaveMode, elts2)

		lpErrsAppender := lpErrorsAppender[T2]{}
		lpErrors = lpErrsAppender.append(lpErrors, failingElts)
		savesCount += uint32(len(elts2) - len(failingElts))
	}

	if t3SaveMode != SaveMode_Noop && len(elts3) > 0 {
		saver := entitiesSaver[T3]{gormDB: streamProcessor.gormDB}
		failingElts := saver.Perform(t3SaveMode, elts3)
		lpErrsAppender := lpErrorsAppender[T3]{}
		lpErrors = lpErrsAppender.append(lpErrors, failingElts)
		savesCount += uint32(len(elts3) - len(failingElts))
	}

	for _, lpError := range lpErrors {
		err := streamProcessor.gormDB.Create(&lpError).Error
		if err != nil {
			log.Printf("Failed to save LineProcessingError %v for %v\n", lpError, err.Error())
		}
	}

	return models.ChunkProcessingReport{SavesCount: savesCount, ErrorsCount: uint32(len(lpErrors))}
}

func (streamProcessor streamProcessorStruct[T1, T2, T3]) Run(
	path string,
	ignoreFirst bool,
	t1SaveMode EntitySaveMode, t2SaveMode EntitySaveMode, t3SaveMode EntitySaveMode,
	reportCh chan models.StreamProcessingReport,
	parseLine ParseLine[T1, T2, T3],
) error {
	var procChunksGathererWG sync.WaitGroup

	chunkProcessingReportsCh := make(chan models.ChunkProcessingReport)

	streamProcessingReport := models.StreamProcessingReport{}
	streamProcessingReport.Path = path
	streamProcessingReport.ChunkProcessingReport = models.ChunkProcessingReport{
		SavesCount:  0,
		ErrorsCount: 0,
	}

	reportCh <- streamProcessingReport

	err := streamProcessor.gormDB.Create(&streamProcessingReport).Error
	if err != nil {
		return err
	}

	go func() {
		for chunkProcessingReport := range chunkProcessingReportsCh {
			streamProcessingReport.SavesCount += chunkProcessingReport.SavesCount
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
						ignoreFirst, chunk, chunkId, parseLine,
						t1SaveMode, t2SaveMode, t3SaveMode,
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
				chunk, chunkId, parseLine,
				t1SaveMode, t2SaveMode, t3SaveMode,
				streamProcessingReport,
			)
		}(chunk, chunkId)
	}

	procChunksGathererWG.Wait()

	if streamProcessingReport.ErrorsCount == 0 {
		streamProcessingReport.Status = streamprocessingstatus.Success
	} else if streamProcessingReport.SavesCount == 0 {
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
