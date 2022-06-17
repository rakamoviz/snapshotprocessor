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
	)

	streamProcessor := MakeStreamProcessor(gormDB, func(path string) (*bufio.Scanner, error) {
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
		streamProcessingReportCh,
		func(line string, gormDB *gorm.DB) error {
			columns := strings.Split(line, ",")
			if len(columns) < 6 {
				return fmt.Errorf("%v: has less than 6 columns", line)
			}

			return gormDB.Transaction(func(tx *gorm.DB) error {
				pCluster := &models.Cluster{
					Code: columns[0][1 : len(columns[0])-1],
				}

				err := gormDB.FirstOrCreate(pCluster, *pCluster).Error
				if err != nil {
					return err
				}

				pNode := &models.Node{
					Code:      columns[1][1 : len(columns[1])-1],
					ClusterID: pCluster.Code,
				}

				err = gormDB.FirstOrCreate(pNode, *pNode).Error
				if err != nil {
					return err
				}

				timestamp, err := strconv.ParseInt(strings.Trim(columns[2], " "), 0, 64)
				if err != nil {
					return err
				}

				pNodeStatus := &models.NodeStatus{
					NodeID: pNode.Code,
					Time:   time.Unix(timestamp, 0),
				}

				err = gormDB.Create(pNodeStatus).Error
				if err != nil {
					return err
				}

				return nil
			})
		},
	)

	for {
		report := <-streamProcessingReportCh
		fmt.Printf("Report. Status: %v, Success:%v, Errors:%v\n", report.Status, report.SuccessCount, report.ErrorsCount)
		if report.Status != streamprocessingstatus.Running && report.Status != streamprocessingstatus.Undefined {
			break
		}
	}
}
