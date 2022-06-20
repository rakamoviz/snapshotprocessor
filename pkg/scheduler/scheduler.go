package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type JobHandler[T any] interface {
	Handle(ctx context.Context, jobData T) error
}

type asynqJobHandler[T any] struct {
	delegate JobHandler[T]
}

func MakeAsynqJobHandler[T any](delegate JobHandler[T]) asynqJobHandler[T] {
	return asynqJobHandler[T]{delegate: delegate}
}

func (h asynqJobHandler[T]) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var jobData T
	if err := json.Unmarshal(t.Payload(), &jobData); err != nil {
		return fmt.Errorf("json.Unmarshal failed %v: %w", err, asynq.SkipRetry)
	}

	return h.delegate.Handle(ctx, jobData)
}

func (h asynqJobHandler[T]) Bind(ctx context.Context, pattern string, s *asynqServer) error {
	s.mux.Handle(pattern, h)
	return nil
}

type Server interface {
	Start(ctx context.Context) error
}

type asynqServer struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func (s *asynqServer) Start(ctx context.Context) error {
	return s.srv.Run(s.mux)
}

func NewAsyncServer(redisClientOpt asynq.RedisClientOpt, config asynq.Config) (*asynqServer, error) {
	srv := asynq.NewServer(redisClientOpt, config)
	mux := asynq.NewServeMux()
	return &asynqServer{srv: srv, mux: mux}, nil
}

type Client[T any] interface {
	Schedule(ctx context.Context, payload T, maxRetry uint8) (string, error)
	Close(ctx context.Context) error
}

type asynqClient[T any] struct {
	pattern string
	client  *asynq.Client
}

func NewAsyncClient[T any](pattern string, redisClientOpt asynq.RedisClientOpt) (*asynqClient[T], error) {
	client := asynq.NewClient(redisClientOpt)
	return &asynqClient[T]{pattern: pattern, client: client}, nil
}

func (c *asynqClient[T]) Schedule(ctx context.Context, jobData T, maxRetry uint8) (string, error) {
	payload, err := json.Marshal(jobData)
	if err != nil {
		return "", err
	}
	task := asynq.NewTask(c.pattern, payload, asynq.MaxRetry(int(maxRetry)))

	info, err := c.client.Enqueue(task)
	if err != nil {
		return "", err
	}

	return info.ID, nil
}

func (c *asynqClient[T]) Close(ctx context.Context) error {
	return c.client.Close()
}
