package main

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/services/streamprocessor"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/workers/streamprocessing"
	"bufio"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"os"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	pkgentities "bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	sqliteDb := sqlite.Open("streamprocessor.tdb?_pragma=busy_timeout(30000)")
	gormDB, err := gorm.Open(sqliteDb, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	gormDB.AutoMigrate(
		&entities.Cluster{}, &entities.Node{}, &entities.NodeStatus{},
		&pkgentities.StreamProcessingReport{}, &pkgentities.LineProcessingError{},
	)
	streamProcessor := streamprocessor.New(gormDB, func(path string) (*bufio.Scanner, error) {
		file, err := os.Open(path)

		if err != nil {
			log.Fatalf("failed to open")

		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		return scanner, nil
	})

	streamProcessingErrorsCh := make(chan error, 10)
	streamProcessingWorker := streamprocessing.New(streamProcessor, 10, streamProcessingErrorsCh)

	go streamProcessingWorker.Start()
	go func() {
		for err := range streamProcessingErrorsCh {
			fmt.Printf("ERROR: %s\n", err.Error())
		}
	}()

	e := echo.New()
	apiGroup := e.Group("/api")
	controllers.Setup(apiGroup, streamProcessingWorker)

	e.Logger.Fatal(e.Start(":1323"))
}
