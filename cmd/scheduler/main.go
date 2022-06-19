package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	internalentities "bitbucket.org/rakamoviz/snapshotprocessor/internal/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/processlines/provider1"
	"bitbucket.org/rakamoviz/snapshotprocessor/internal/scheduler/handlers"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/scheduler"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/services/streamprocessor"
	"github.com/glebarez/sqlite"
	"github.com/hibiken/asynq"
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

	streamProcessor := streamprocessor.New(gormDB, func(path string) (*bufio.Scanner, error) {
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
		handlers.NewStreamProcessing(streamProcessor, processLines),
	)
	streamProcessingJobHandler.Bind(string(handlers.StreamProcessing), asynqServer)

	fmt.Println(">>>> Scheduler starting...")
	err = asynqServer.Start()
	if err != nil {
		log.Fatalln(err.Error())
	}
}
