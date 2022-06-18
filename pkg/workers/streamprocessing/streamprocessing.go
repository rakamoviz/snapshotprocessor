package streamprocessing

import (
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/entities"
	"bitbucket.org/rakamoviz/snapshotprocessor/pkg/services/streamprocessor"
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
	for job := range worker.jobsCh {
		err := worker.sproc.Run(job.path, job.ignoreFirst, job.reportsCh, job.processLine)
		if err != nil {
			worker.errorsCh <- err
		}
	}
}

func (worker *worker) AppendJob(path string, ignoreFirst bool, processLine streamprocessor.ProcessLine) <-chan entities.StreamProcessingReport {
	reportsCh := make(chan entities.StreamProcessingReport)
	worker.jobsCh <- job{path: path, ignoreFirst: ignoreFirst, reportsCh: reportsCh, processLine: processLine}

	return reportsCh
}
