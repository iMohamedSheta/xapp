package bus

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
)

const defaultTaskPrefix = "bus:"

// QueueBackend interface should be implemented by the queue provider (e.g. asynq, redis, rabbitmq).
type QueueBackend interface {
	Enqueue(ctx context.Context, listener Listener, taskKey string, payload []byte) error
}

// Listener interface should be implemented by the listeners.
type Listener interface {
	Handle(ctx context.Context, payload []byte) error
	ShouldQueue() bool
	TaskName() string
}

// Bus manages event subscriptions and dispatch.
type Bus struct {
	mu        sync.RWMutex
	listeners map[string][]Listener
	byTask    map[string]Listener
	queue     QueueBackend
	prefix    string
}

// TaskListeners returns a snapshot of every queueing listener, keyed by
// TaskName. Backends that route by queue-native task type — registering
// one handler per TaskName instead of going through HandleTask — use this
// at boot to build their registration table.
func (b *Bus) TaskListeners() map[string]Listener {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make(map[string]Listener, len(b.byTask))
	for name, l := range b.byTask {
		out[name] = l
	}
	return out
}

// New creates a new Bus instance.
func New(queue QueueBackend, prefix string) *Bus {
	if prefix == "" {
		prefix = defaultTaskPrefix
	}
	return &Bus{
		listeners: make(map[string][]Listener),
		byTask:    make(map[string]Listener),
		queue:     queue,
		prefix:    prefix,
	}
}

// Subscribe subscribes a listener to an event.
// If the listener ShouldQueue() is true it will be queued.
// Otherwise, it will be handled synchronously.
func (b *Bus) Subscribe(event string, l Listener) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if l.ShouldQueue() {
		name := l.TaskName()
		if name == "" {
			return fmt.Errorf("bus: listener %T has ShouldQueue=true but empty TaskName", l)
		}
		key := b.taskKey(name)
		if existing, ok := b.byTask[key]; ok && existing != l {
			return fmt.Errorf("bus: duplicate TaskName %q (existing %T, new %T)", key, existing, l)
		}
		b.byTask[key] = l
	}
	b.listeners[event] = append(b.listeners[event], l)
	return nil
}

// Publish dispatches an event to all listeners subscribed to it.
// If the listener ShouldQueue() is true it will be queued.
// Otherwise, it will be handled synchronously.
func (b *Bus) Publish(ctx context.Context, event string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	b.mu.RLock()
	ls := b.listeners[event]
	b.mu.RUnlock()

	var errs []error
	for _, l := range ls {
		if l.ShouldQueue() {
			if b.queue == nil {
				errs = append(errs, fmt.Errorf("bus: no queue backend, dropping %q", l.TaskName()))
				continue
			}
			if err := b.queue.Enqueue(ctx, l, b.taskKey(l.TaskName()), payload); err != nil {
				errs = append(errs, err)
			}
			continue
		}
		if err := l.Handle(ctx, payload); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// HandleTask gets the listener by task name and handles the event.
func (b *Bus) HandleTask(ctx context.Context, taskName string, payload []byte) error {
	b.mu.RLock()
	l, ok := b.byTask[taskName]
	b.mu.RUnlock()
	if !ok {
		return fmt.Errorf("bus: no listener for task %q", taskName)
	}
	return l.Handle(ctx, payload)
}

// taskKey returns the canonical queue task type for a TaskName, scoped to
// this Bus's own prefix. Subscribe and Publish are the only two callers —
// backends never compute this themselves.
func (b *Bus) taskKey(name string) string {
	if strings.HasPrefix(name, b.prefix) {
		return name
	}
	return b.prefix + name
}

// Prefix exposes the task-key prefix this Bus was constructed with, so
// backends can register a queue-native handler under a matching pattern.
func (b *Bus) Prefix() string {
	return b.prefix
}

/*
	the api will be like this

	// Create the queue backend for the bus
	queueClient, err := xioc.AppMake[*x.QueueClient]()
	backend := bus.NewAsynqBackend(queueClient.Client, asynq.Queue(tasks.QueueEvents), asynq.MaxRetry(3))
    // register the bus inside the ioc and use it like x.Bus()
	bus.New(backend)

	// this is how we dispatch an event happened,
	// this is sync
	x.Bus().Dispatch(ctx, events.EventSomethingHappend, data as any)
	// this is async (queue)
	x.Bus().Publish(ctx, events.EventSomethingHappend, data as any)

	but i already have a ShouldQueue() so i don't need another method like dispatch

	so lets say now we dispatched the event with the payload (data)
	now we need handler to handle it the handler should implement the Listener interface

	type SomethingHappendListener struct{
		authService *AuthService
		auditService *auditService
		userRepo    *userRepo
	}

	func (l *SomethingHappendListener) Handle(ctx context.Context, payload []byte) error {
		var data struct{
			UserId int64

		}
		return nil
	}

	func (l *SomethingHappendListener) TaskName() string {
		return "something_happend"
	}

	func (l *SomethingHappendListener) ShouldQueue() bool {
		return true
	}

	// optional to override the default options
	func (l *SomethingHappendListener) AsynqOpts() []asynq.Option {
		return []asynq.Option{
			asynq.Queue(tasks.QueueEvents),
			asynq.MaxRetry(3),
		}
	}

	and we should register this Event listener inside the registers file

	// BusListeners subscribes every domain event listener to the bus.
	// Call once after InitBus has run.
	func BusListeners() {
		b := x.Bus()

		subscribe := func(event string, l Listener) {
			if err := b.Subscribe(event, l); err != nil {
				errMsg := "bus: subscribe failed: " + err.Error()
				x.Logger().Error(errMsg)
				utils.PrintErr(errMsg)
			}
		}

		subscribe(events.EventSomethingHappend, NewSomethingHappendListener() || xioc.AppMake[*SomethingHappendListener]())
	}

	and now when the bus try to handle it will give the handler the data which is Published with it as handler signature
	Handle(ctx context.Context, payload []byte)

*/
