package main

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/streamprocessor"
	"bufio"
	"github.com/labstack/echo/v4"
	"log"
	"os"

	"bitbucket.org/rakamoviz/snapshotprocessor/internal/db/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	sqliteDb := sqlite.Open("test_streamprocessor.tdb?_pragma=busy_timeout(30000)")
	gormDB, err := gorm.Open(sqliteDb, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	gormDB.AutoMigrate(
		&models.Cluster{}, &models.Node{}, &models.NodeStatus{},
		&models.StreamProcessingReport{}, &models.LineProcessingError{},
	)
	streamProcessor := streamprocessor.MakeStreamProcessor(gormDB, func(path string) (*bufio.Scanner, error) {
		file, err := os.Open(path)

		if err != nil {
			log.Fatalf("failed to open")

		}
		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		return scanner, nil
	})

	e := echo.New()
	apiGroup := e.Group("/api")
	controllers.Setup(apiGroup, streamProcessor)

	e.Logger.Fatal(e.Start(":1323"))
}
