package main

import (
	"os"

	"bitbucket.org/rakamoviz/snapshotprocessor/cmd/middleware/controllers"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/scheduler"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
)

func main() {
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		redisAddr = "127.0.0.1:6379"
	}

	streamProcessingScheduler, _ := scheduler.NewAsyncClient[handlers.StreamProcessingJobData](
		string(handlers.StreamProcessing),
		asynq.RedisClientOpt{
			Addr: redisAddr,
		},
	)

	e := echo.New()
	apiGroup := e.Group("/api")
	controllers.Setup(apiGroup, streamProcessingScheduler)

	e.Logger.Fatal(e.Start(":1323"))
}
