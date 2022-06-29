package main

import (
	"bufio"
	"log"
	"os"

	"context"

	"github.com/glebarez/sqlite"
	"github.com/hibiken/asynq"
	internalentities "github.com/rakamoviz/snapshotprocessor/internal/entities"
	"github.com/rakamoviz/snapshotprocessor/internal/processlines/provider1"
	internalhandlers "github.com/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	"github.com/rakamoviz/snapshotprocessor/pkg/entities"
	"github.com/rakamoviz/snapshotprocessor/pkg/scheduler"
	"github.com/rakamoviz/snapshotprocessor/pkg/scheduler/handlers"
	"github.com/rakamoviz/snapshotprocessor/pkg/services/streamprocessor"
	"gorm.io/gorm"
)

func main() {
	redisAddr, ok := os.LookupEnv("REDIS_ADDR")
	if !ok {
		redisAddr = "127.0.0.1:6379"
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

	streamProcessor := streamprocessor.New(gormDB, func(ctx context.Context, path string) (*bufio.Scanner, error) {
		file, err := os.Open(path)

		if err != nil {
			//log.Fatalf("failed to open", err, path)
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)
		return scanner, nil
	})

	asynqServer, err := scheduler.NewAsyncServer(asynq.RedisClientOpt{
		Addr: redisAddr,
	}, asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create AsyncServer %v", err)
	}

	processLines := map[string]map[string]streamprocessor.ProcessLine{
		"provider1": {
			"format1": provider1.ProcessFormat1,
		},
	}

	streamProcessingJobHandler := scheduler.MakeAsynqJobHandler[handlers.StreamProcessingJobData](
		handlers.NewStreamProcessing(gormDB, streamProcessor, processLines),
	)

	ctx := context.Background()
	streamProcessingJobHandler.Bind(ctx, string(internalhandlers.StreamProcessing), asynqServer)

	err = asynqServer.Start(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
