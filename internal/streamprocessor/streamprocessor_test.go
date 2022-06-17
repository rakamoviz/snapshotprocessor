package streamprocessor

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models/streamprocessingstatus"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestStreamProcessor(t *testing.T) {
	sqliteDb := sqlite.Open("test_streamprocessor.tdb?_pragma=busy_timeout(30000)")
	gormDB, err := gorm.Open(sqliteDb, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	gormDB.AutoMigrate(
		&models.Cluster{}, &models.Node{}, &models.NodeStatus{},
		&models.StreamProcessingReport{}, &models.LineProcessingError{},
		&models.ProcessedLine{},
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
		"/home/rcokorda/Projects/snapshotprocessor/sandbox/snapshots.csv",
		true,
		SaveMode_InsertIfInexist, SaveMode_InsertIfInexist, SaveMode_Insert,
		streamProcessingReportCh,
		func(line string) (*models.Cluster, *models.Node, *models.NodeStatus, error) {
			columns := strings.Split(line, ",")
			if len(columns) < 6 {
				return nil, nil, nil, fmt.Errorf("%v: has less than 6 columns", line)
			}

			pCluster := &models.Cluster{
				Code: columns[0][1 : len(columns[0])-1],
			}

			pNode := &models.Node{
				Code:      columns[1][1 : len(columns[1])-1],
				ClusterID: pCluster.Code,
			}

			timestamp, err := strconv.ParseInt(strings.Trim(columns[2], " "), 0, 64)
			if err != nil {
				return nil, nil, nil, err
			}

			pNodeStatus := &models.NodeStatus{
				NodeID: pNode.Code,
				Time:   time.Unix(timestamp, 0),
			}

			return pCluster, pNode, pNodeStatus, nil
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
