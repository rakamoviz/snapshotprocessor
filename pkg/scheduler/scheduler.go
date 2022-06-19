package scheduler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

type JobHandler[T any] interface {
	Handle(jobData T) error
}

type asynqJobHandler[T any] struct {
	delegate JobHandler[T]
}

func MakeAsynqJobHandler[T any](delegate JobHandler[T]) asynqJobHandler[T] {
	return asynqJobHandler[T]{delegate: delegate}
}

func (h asynqJobHandler[T]) ProcessTask(ctx context.Context, t *asynq.Task) error {
	fmt.Println("sshsishsshsihsishis ")
	var jobData T
	if err := json.Unmarshal(t.Payload(), &jobData); err != nil {
		fmt.Println("ABABABABAB ", err)
		return fmt.Errorf("json.Unmarshal failed %v: %w", err, asynq.SkipRetry)
	}

	fmt.Println("XXXXXXXXXXXXXXx ", h.delegate, jobData)
	return h.delegate.Handle(jobData)
}

func (h asynqJobHandler[T]) Bind(pattern string, s *asynqServer) error {
	fmt.Println("!!!!!!!!!!!!!!!!!!! ", pattern, h)
	s.mux.Handle(pattern, h)
	return nil
}

type Server interface {
	Start() error
}

type asynqServer struct {
	srv *asynq.Server
	mux *asynq.ServeMux
}

func (s *asynqServer) Start() error {
	return s.srv.Run(s.mux)
}

func NewAsyncServer(redisClientOpt asynq.RedisClientOpt, config asynq.Config) (*asynqServer, error) {
	srv := asynq.NewServer(redisClientOpt, config)
	mux := asynq.NewServeMux()
	return &asynqServer{srv: srv, mux: mux}, nil
}

type Client[T any] interface {
	Schedule(payload T) (string, error)
	Close() error
}

type asynqClient[T any] struct {
	pattern string
	client  *asynq.Client
}

func NewAsyncClient[T any](pattern string, redisClientOpt asynq.RedisClientOpt) (*asynqClient[T], error) {
	client := asynq.NewClient(redisClientOpt)
	return &asynqClient[T]{pattern: pattern, client: client}, nil
}

func (c *asynqClient[T]) Schedule(jobData T) (string, error) {
	payload, err := json.Marshal(jobData)
	if err != nil {
		return "", err
	}
	task := asynq.NewTask(c.pattern, payload)

	info, err := c.client.Enqueue(task)
	fmt.Println(">>>>>>>>>>>>>>>>>>>> ", task, c.pattern)
	if err != nil {
		return "", err
	}

	return info.ID, nil
}

func (c *asynqClient[T]) Close() error {
	return c.client.Close()
}