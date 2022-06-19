package main

import (
	"os"

	"log"

	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers"
	internalentities "bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/scheduler"
	"github.com/glebarez/sqlite"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func main() {
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		redisAddr = "127.0.0.1:6379"
	}

	streamProcessingScheduler, err := scheduler.NewAsyncClient[handlers.StreamProcessingJobData](
		string(handlers.StreamProcessing),
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	sqliteDb := sqlite.Open("/home/rcokorda/Projects/snapshotprocessor/snapshotprocessor.tdb?_pragma=busy_timeout(30000)")
	gormDB, err := gorm.Open(sqliteDb, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	gormDB.AutoMigrate(
		&internalentities.Cluster{}, &internalentities.Node{}, &internalentities.NodeStatus{},
		&entities.StreamProcessingReport{}, &entities.LineProcessingError{},
	)

	e := echo.New()
	apiGroup := e.Group("/api")
	controllers.Setup(apiGroup, gormDB, streamProcessingScheduler)

	e.Logger.Fatal(e.Start(":1323"))
}
