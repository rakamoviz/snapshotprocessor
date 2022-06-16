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

type StreamProcessor[T1 any, T2 any, T3 any] interface {
	Run(
		path string,
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

func (saver entitiesSaver[T]) Perform(saveMode EntitySaveMode, entities []*T) (uint32, uint32) {
	savesCount := uint32(0)
	errorsCount := uint32(0)

	switch saveMode {
	case SaveMode_Insert:
		for _, entity := range entities {
			tx := saver.gormDB.Create(entity)
			if tx.Error == nil {
				savesCount++
			} else {
				errorsCount++
			}
		}
	case SaveMode_InsertIfInexist:
		for _, entity := range entities {
			tx := saver.gormDB.FirstOrCreate(entity)
			if tx.Error == nil {
				savesCount++
			} else {
				errorsCount++
			}
		}
	case SaveMode_Update:
		for _, entity := range entities {
			tx := saver.gormDB.Save(entity)
			if tx.Error == nil {
				savesCount++
			} else {
				errorsCount++
			}
		}
	case SaveMode_Upsert:
		for _, entity := range entities {
			tx := saver.gormDB.Save(entity)
			if tx.Error == nil {
				savesCount++
			} else {
				errorsCount++
			}
		}
	}
	return savesCount, errorsCount
}

func MakeStreamProcessor[T1 any, T2 any, T3 any](gormDB *gorm.DB, openScanner OpenScanner) StreamProcessor[T1, T2, T3] {
	return streamProcessorStruct[T1, T2, T3]{gormDB: gormDB, openScanner: openScanner}
}

func (streamProcessor streamProcessorStruct[T1, T2, T3]) processChunk(
	chunk []string, chunkId int, parseLine ParseLine[T1, T2, T3],
	t1SaveMode EntitySaveMode, t2SaveMode EntitySaveMode, t3SaveMode EntitySaveMode,
	report models.StreamProcessingReport,
) models.ChunkProcessingReport {
	lineNumberOffset := chunkId * CHUNK_LEN

	var entities1 []*T1
	var entities2 []*T2
	var entities3 []*T3

	if t1SaveMode != SaveMode_Noop {
		entities1 = make([]*T1, len(chunk))
	}

	if t2SaveMode != SaveMode_Noop {
		entities2 = make([]*T2, len(chunk))
	}

	if t3SaveMode != SaveMode_Noop {
		entities3 = make([]*T3, len(chunk))
	}

	lpErrors := make([]models.LineProcessingError, len(chunk))

	lpErrorsIndex := 0
	entities1Index := 0
	entities2Index := 0
	entities3Index := 0

	for i, line := range chunk {
		entity1, entity2, entity3, err := parseLine(line)
		if err != nil {
			lpErrors[lpErrorsIndex] = models.LineProcessingError{
				LineNumber: lineNumberOffset + i,
				Line:       line,
				Error:      err.Error(),
				Report:     report,
			}
			lpErrorsIndex++
		} else {
			if t1SaveMode != SaveMode_Noop && entity1 != nil {
				entities1[entities1Index] = entity1
				entities1Index++
			}

			if t2SaveMode != SaveMode_Noop && entity2 != nil {
				entities2[entities2Index] = entity2
				entities2Index++
			}

			if t3SaveMode != SaveMode_Noop && entity3 != nil {
				entities3[entities3Index] = entity3
				entities3Index++
			}
		}
	}

	//todo: update SuccessfulInsertsCount and ErrorsCount in the database

	savesCount := uint32(0)
	errorsCount := uint32(0)

	if t1SaveMode != SaveMode_Noop && entities1Index > 0 {
		saver := entitiesSaver[T1]{gormDB: streamProcessor.gormDB}
		sc, ec := saver.Perform(t1SaveMode, entities1[0:entities1Index:len(chunk)])

		savesCount += sc
		errorsCount += ec
	}

	if t2SaveMode != SaveMode_Noop && entities2Index > 0 {
		saver := entitiesSaver[T2]{gormDB: streamProcessor.gormDB}
		sc, ec := saver.Perform(t2SaveMode, entities2[0:entities2Index:len(chunk)])

		savesCount += sc
		errorsCount += ec
	}

	if t3SaveMode != SaveMode_Noop && entities3Index > 0 {
		saver := entitiesSaver[T3]{gormDB: streamProcessor.gormDB}
		sc, ec := saver.Perform(t3SaveMode, entities3[0:entities3Index:len(chunk)])

		savesCount += sc
		errorsCount += ec
	}

	if lpErrorsIndex > 0 {
		errLpErrors := streamProcessor.gormDB.Create(lpErrors[0:lpErrorsIndex:len(chunk)]).Error
		if errLpErrors != nil {
			log.Printf("Failed saving LineProcessingErrors\n") //TODO: provide more detail
		}
	}

	//finalErrorsCount := errorsCount + (insertsCount - successfulInsertsCount)
	return models.ChunkProcessingReport{SavesCount: savesCount, ErrorsCount: errorsCount}
}

func (streamProcessor streamProcessorStruct[T1, T2, T3]) Run(
	path string,
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
	chunkId := -1
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
		line := scanner.Text()

		if chunkLineOffset%CHUNK_LEN == 0 {
			if chunkId >= 0 {
				procChunksGathererWG.Add(1)
				go func(chunk []string, chunkId int) {
					chunkProcessingReportsCh <- streamProcessor.processChunk(
						chunk, chunkId, parseLine,
						t1SaveMode, t2SaveMode, t3SaveMode,
						streamProcessingReport,
					)
				}(chunk, chunkId)
			}

			chunkLineOffset = 0
			chunk = make([]string, CHUNK_LEN)
			chunkId += 1
		}

		chunk[chunkLineOffset] = line
		chunkLineOffset += 1
	}

	if chunkId >= 0 {
		procChunksGathererWG.Add(1)
		go func(chunk []string, chunkId int) {
			chunkProcessingReportsCh <- streamProcessor.processChunk(
				chunk[0:chunkLineOffset:CHUNK_LEN], chunkId, parseLine,
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
