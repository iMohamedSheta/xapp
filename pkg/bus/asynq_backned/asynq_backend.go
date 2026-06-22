package asynq_backend

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/imohamedsheta/xapp/pkg/bus"
)

type AsynqOpts interface {
	AsynqOpts() []asynq.Option
}

type Enqueuer interface {
	EnqueueContext(ctx context.Context, task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
}

type AsynqBackend struct {
	client      Enqueuer
	defaultOpts []asynq.Option
}

func NewAsynqBackend(client Enqueuer, defaultOpts ...asynq.Option) *AsynqBackend {
	return &AsynqBackend{client: client, defaultOpts: defaultOpts}
}

// func (b *AsynqBackend) Enqueue(ctx context.Context, listener bus.Listener, payload []byte) error {
// 	data, err := json.Marshal(bus.BusTask{TaskName: listener.TaskName(), Payload: payload})
// 	if err != nil {
// 		return fmt.Errorf("bus - asynq: %w", err)
// 	}

// 	opts := append([]asynq.Option{}, b.defaultOpts...)
// 	if a, ok := listener.(AsynqOpts); ok {
// 		opts = append(opts, a.AsynqOpts()...)
// 	}

// 	if _, err := b.client.EnqueueContext(ctx, asynq.NewTask(bus.TaskType, data), opts...); err != nil {
// 		return fmt.Errorf("bus - asynq: enqueue: %w", err)
// 	}
// 	return nil
// }

// Enqueue submits the job under the listener's own TaskName as the asynq
// task type — no wrapper struct needed, the type itself carries the routing
// that BusTask used to carry inside the payload.
func (b *AsynqBackend) Enqueue(ctx context.Context, listener bus.Listener, taskKey string, payload []byte) error {
	opts := append([]asynq.Option{}, b.defaultOpts...)
	if a, ok := listener.(AsynqOpts); ok {
		opts = append(opts, a.AsynqOpts()...)
	}
	if _, err := b.client.EnqueueContext(ctx, asynq.NewTask(taskKey, payload), opts...); err != nil {
		return fmt.Errorf("bus - asynq: enqueue: %w", err)
	}
	return nil
}

// AsynqHandler is the *only* thing in this package that implements
// asynq.Handler. Bus itself stays free of any asynq-shaped method —
// swap backends and this is the one type that gets replaced, not Bus.
type AsynqHandler struct {
	bus *bus.Bus
}

func NewAsynqHandler(b *bus.Bus) *AsynqHandler {
	return &AsynqHandler{bus: b}
}

func (h *AsynqHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	return h.bus.HandleTask(ctx, t.Type(), t.Payload())
}

// func (h *AsynqHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
// 	var task bus.BusTask
// 	if err := json.Unmarshal(t.Payload(), &task); err != nil {
// 		return fmt.Errorf("bus - asynq: %w", err)
// 	}
// 	return h.bus.HandleTask(ctx, task.TaskName, task.Payload)
// }
