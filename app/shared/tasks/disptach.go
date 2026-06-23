package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/imohamedsheta/xapp/pkg/chainq"
)

// // Dispatch a new task to the queue
// func DispatchAsynqTask(task *asynq.Task, opts ...asynq.Option) error {
// 	client := x.Queue().Client
// 	if client == nil {
// 		return fmt.Errorf("asynq client not initialized")
// 	}

// 	info, err := client.Enqueue(task, opts...)
// 	if err != nil {
// 		return fmt.Errorf("failed to enqueue task: %w", err)
// 	}

// 	log.Printf("Task enqueued: ID=%s, Queue=%s, Type=%s",
// 		info.ID, info.Queue, info.Type)
// 	return nil
// }

/*
	Helpers to create new asynq Task or TaskHandler used inside the tasks
*/

// Create new asynq task to from task to use it inside package
func CreateAsynqTask(t chainq.Task, opts ...asynq.Option) (*asynq.Task, error) {
	payload, err := json.Marshal(t.GetPayload())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	return asynq.NewTask(t.GetTaskType(), payload, opts...), nil
}

// Process task payload helper to send payload to actual handler
func processTaskPayload[T any](ctx context.Context, task *asynq.Task, handler func(context.Context, *T) error) error {
	var payload T
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return handler(ctx, &payload)
}
