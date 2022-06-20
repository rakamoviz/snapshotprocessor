package streamprocessor

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	internalentities "bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities/streamprocessingstatus"
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
		&internalentities.Cluster{}, &internalentities.Node{}, &internalentities.NodeStatus{},
		&entities.StreamProcessingReport{}, &entities.LineProcessingError{},
	)

	streamProcessor := New(gormDB, func(ctx context.Context, path string) (*bufio.Scanner, error) {
		file, err := os.Open(path)

		if err != nil {
			log.Fatalf("failed to open")

		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		return scanner, nil
	})

	streamProcessingReportCh := make(chan entities.StreamProcessingReport)
	errorsCh := make(chan error)
	go streamProcessor.Run(
		context.Background(),
		"/home/rcokorda/Projects/snapshotprocessor/sandbox/snapshots.csv",
		true,
		streamProcessingReportCh,
		errorsCh,
		func(line string, gormDB *gorm.DB) error {
			columns := strings.Split(line, ",")
			if len(columns) < 6 {
				return fmt.Errorf("%v: has less than 6 columns", line)
			}

			return gormDB.Transaction(func(tx *gorm.DB) error {
				cluster := internalentities.Cluster{
					Code: columns[0][1 : len(columns[0])-1],
				}

				err := gormDB.FirstOrCreate(&cluster, cluster).Error
				if err != nil {
					return err
				}

				node := internalentities.Node{
					Code:      columns[1][1 : len(columns[1])-1],
					ClusterID: cluster.Code,
				}

				err = gormDB.FirstOrCreate(&node, node).Error
				if err != nil {
					return err
				}

				timestamp, err := strconv.ParseInt(strings.Trim(columns[2], " "), 0, 64)
				if err != nil {
					return err
				}

				nodeStatus := internalentities.NodeStatus{
					NodeID: node.Code,
					Time:   time.Unix(timestamp, 0),
				}

				err = gormDB.Create(&nodeStatus).Error
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
