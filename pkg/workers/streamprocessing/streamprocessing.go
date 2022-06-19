package streamprocessing

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/services/streamprocessor"
	"fmt"
)

type job struct {
	path        string
	ignoreFirst bool
	processLine streamprocessor.ProcessLine
	reportsCh   chan entities.StreamProcessingReport
}

type Worker interface {
	Start()
	AppendJob(path string, ignoreFirst bool, processLine streamprocessor.ProcessLine) <-chan entities.StreamProcessingReport
}

type worker struct {
	sproc    streamprocessor.StreamProcessor
	jobsCh   chan job
	errorsCh chan<- error
}

func New(
	sproc streamprocessor.StreamProcessor,
	size uint8, errorsCh chan<- error,
) Worker {
	return &worker{
		sproc:    sproc,
		jobsCh:   make(chan job, size),
		errorsCh: errorsCh,
	}
}

func (worker *worker) Start() {
	fmt.Println("c.1")
	for job := range worker.jobsCh {
		fmt.Println("c.2")
		go worker.sproc.Run(job.path, job.ignoreFirst, job.reportsCh, worker.errorsCh, job.processLine)
		fmt.Println("c.3")
	}
}

func (worker *worker) AppendJob(path string, ignoreFirst bool, processLine streamprocessor.ProcessLine) <-chan entities.StreamProcessingReport {
	reportsCh := make(chan entities.StreamProcessingReport)
	fmt.Println("A.1")
	worker.jobsCh <- job{path: path, ignoreFirst: ignoreFirst, reportsCh: reportsCh, processLine: processLine}
	fmt.Println("A.2")

	return reportsCh
}
