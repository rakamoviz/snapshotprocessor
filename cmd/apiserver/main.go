package main

import (
	"os"

	"log"

	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/apiserver/controllers"
	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/apiserver/middlewares"
	internalentities "bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	internalhandlers "bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	pkgentities "bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/repository"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/scheduler"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/scheduler/handlers"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/services/auth"
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
		string(internalhandlers.StreamProcessing),
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
		&pkgentities.StreamProcessingReport{}, &pkgentities.LineProcessingError{},
	)

	apiKeyCheck := middlewares.NewApiKeyCheck(auth.NewMemoryBasedService(map[string]auth.ApiClient{
		"abcdef": {Name: "provider1"},
	}))

	clusterRepository := repository.New[internalentities.Cluster](gormDB)
	nodeRepository := repository.New[internalentities.Node](gormDB)
	nodeStatusRepository := repository.New[internalentities.NodeStatus](gormDB)
	streamProcessingReportRepository := repository.New[pkgentities.StreamProcessingReport](gormDB)
	lineProcessingErrorRepository := repository.New[pkgentities.LineProcessingError](gormDB)

	e := echo.New()
	apiGroup := e.Group("/api")
	controllers.Setup(
		apiGroup, gormDB, streamProcessingScheduler, apiKeyCheck,
		clusterRepository, nodeRepository, nodeStatusRepository, streamProcessingReportRepository,
		lineProcessingErrorRepository,
	)

	e.Logger.Fatal(e.Start(":1323"))
}
