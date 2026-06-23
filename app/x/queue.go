package x

import (
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/imohamedsheta/xapp/app/shared/adapters"
	"github.com/imohamedsheta/xapp/pkg/chainq"
)

/*
|--------------------------------------------------------
|	Application Dependency Container Alias
|--------------------------------------------------------
*/

type QueueClient struct {
	Client *asynq.Client
}

/*
|--------------------------------------------------------
|	Application Dependency Container Calls
|--------------------------------------------------------
|*/

func NewQueueClient(client *asynq.Client) *QueueClient {
	return &QueueClient{Client: client}
}

func Queue() *QueueClient {
	queue, err := app[*QueueClient]()
	if err != nil {
		Logger().Error(fmt.Sprintf("QueueClient can't be resolved: %s", err.Error()))
		return nil
	}
	return queue
}

// This implements the Shutdownable interface the ioc shutdown the service when the application is shutdown
func (c *QueueClient) Shutdown() error {
	if c.Client != nil {
		return c.Client.Close()
	}
	return nil
}

/*
|--------------------------------------------------------
|  Dispatchers Helpers
|--------------------------------------------------------
*/

// Dispatch a new task to the queue
func Dispatch(task chainq.Task, opts ...asynq.Option) error {
	client := Queue().Client
	if client == nil {
		return fmt.Errorf("asynq client not initialized")
	}

	asynqTask, err := task.CreateTask()
	if err != nil {
		return err
	}

	info, err := client.Enqueue(asynqTask, opts...)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("Task enqueued: ID=%s, Queue=%s, Type=%s",
		info.ID, info.Queue, info.Type)
	return nil
}

// Chain - helper to create new chainq.Chain which can be used to create new chain of tasks
func Chain(opt *chainq.ChainOptions) *chainq.Chain {
	if opt == nil {
		opt = &chainq.ChainOptions{
			MaxRetries:   3,
			Timeout:      5 * time.Second,
			DefaultQueue: "default",
		}
	}

	return chainq.NewChain(
		Queue().Client,
		adapters.NewLoggerAdapter(Logger()),
		opt,
	)
}
