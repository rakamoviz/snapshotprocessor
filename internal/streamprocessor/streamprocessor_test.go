package streamprocessor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"testing"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models/streamprocessingstatus"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestStreamProcessor(t *testing.T) {
	sqliteDb := sqlite.Open("test_streamprocessor.tdb")
	gormDB, err := gorm.Open(sqliteDb, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	gormDB.AutoMigrate(
		&models.Cluster{}, &models.Node{}, &models.NodeStatus{},
		&models.StreamProcessingReport{}, &models.LineProcessingError{},
	)

	streamProcessor := MakeStreamProcessor[models.Cluster, models.Node, models.NodeStatus](gormDB, func(path string) (*bufio.Scanner, error) {
		file, err := os.Open(path)

		if err != nil {
			log.Fatalf("failed to open")

		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		return scanner, nil
	})

	streamProcessingReportCh := make(chan models.StreamProcessingReport)
	go streamProcessor.Run(
		"/home/rcokorda/Projects/snapshotprocessor/sandbox/dataset.csv",
		SaveMode_InsertIfInexist, SaveMode_InsertIfInexist, SaveMode_Insert,
		streamProcessingReportCh,
		func(line string) (*models.Cluster, *models.Node, *models.NodeStatus, error) {
			return nil, nil, nil, nil
		},
	)

	for {
		report := <-streamProcessingReportCh
		fmt.Printf("Report. Status: %v, Saves:%v, Errors:%v\n", report.Status, report.SavesCount, report.ErrorsCount)
		if report.Status != streamprocessingstatus.Running && report.Status != streamprocessingstatus.Undefined {
			break
		}
	}
}
