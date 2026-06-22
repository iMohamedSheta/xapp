package goroutine_backend

import (
	"context"
	"log/slog"
	"sync"

	"github.com/imohamedsheta/xapp/pkg/bus"
)

// Logger is a minimal structured logger interface.
// Compatible with slog.Logger, zap.SugaredLogger, and similar.
type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Warn(msg string, fields ...any)
}

// QueueHandler is satisfied by *bus.Bus via HandleTask. Kept narrow (not
// *bus.Bus directly) for the same reason Enqueuer exists in asynq_backend:
// swap or mock the dispatch side without depending on the concrete type.
// Named QueueHandler, not Handler, because engine.go already owns Handler[T]
// at package scope.
type QueueHandler interface {
	HandleTask(ctx context.Context, taskKey string, payload []byte) error
}

// GoroutineOptions configures the GoroutineBackend. Zero-value Workers/
// QueueSize fall through to xqueue's own defaults, so no need to
// redeclare DefaultWorkers/DefaultQueueSize in this package.
type GoroutineOptions struct {
	Workers   int
	QueueSize int
	Logger    Logger
}

func (o GoroutineOptions) logger() Logger {
	if o.Logger == nil {
		return slog.Default()
	}
	return o.Logger
}

type goroutineItem struct {
	listener bus.Listener
	taskKey  string
	payload  []byte
}

// GoroutineBackend implements QueueBackend using an in-process xqueue.Engine.
// Call SetHandler once after construction, before any Enqueue. Call
// Shutdown to drain and stop gracefully.
type GoroutineBackend struct {
	mu      sync.RWMutex
	engine  *Engine[goroutineItem]
	handler QueueHandler
	logger  Logger
}

// NewGoroutineBackend creates a GoroutineBackend and launches its worker pool.
// The returned backend has no handler yet — call SetHandler before
// publishing anything, or every dequeued item will fail with
// "no handler registered".
func NewGoroutineBackend(opts GoroutineOptions) *GoroutineBackend {
	b := &GoroutineBackend{logger: opts.logger()}
	b.engine = New(Options{
		Workers:   opts.Workers,
		QueueSize: opts.QueueSize,
		Logger:    b.logger,
	}, b.process)
	return b
}

func (b *GoroutineBackend) process(ctx context.Context, it goroutineItem) error {
	b.mu.RLock()
	h := b.handler
	b.mu.RUnlock()

	if h != nil {
		if err := h.HandleTask(ctx, it.taskKey, it.payload); err != nil {
			b.logger.Error("bus - goroutine backend: handler error", "task_key", it.taskKey, "error", err)
			return err
		}
		return nil
	}

	// Fallback: no handler wired, dispatch straight to the listener like
	// before HandleTask existed. Logged at Warn so a forgotten SetHandler
	// call shows up somewhere instead of disappearing entirely.
	b.logger.Warn("bus - goroutine backend: no handler registered, dispatching directly",
		"task_name", it.listener.TaskName(),
	)
	if err := it.listener.Handle(ctx, it.payload); err != nil {
		b.logger.Error("bus - goroutine backend: handler error", "task_name", it.listener.TaskName(), "error", err)
		return err
	}
	return nil
}

// SetHandler wires the dispatch side. Call once, right after bus.New(b, prefix)
// returns, before anything is published. Panics on nil — a backend with no
// handler is a wiring bug, not a runtime condition.
func (b *GoroutineBackend) SetHandler(h QueueHandler) {
	if h == nil {
		panic("bus - goroutine backend: SetHandler called with nil handler")
	}
	b.mu.Lock()
	b.handler = h
	b.mu.Unlock()
}

// Enqueue implements QueueBackend.
func (b *GoroutineBackend) Enqueue(ctx context.Context, listener bus.Listener, taskKey string, payload []byte) error {
	return b.engine.Enqueue(ctx, goroutineItem{listener: listener, taskKey: taskKey, payload: payload}, nil)
}

// Shutdown waits for in-flight tasks to finish, or returns ctx.Err() if the
// context is cancelled first.
func (b *GoroutineBackend) Shutdown(ctx context.Context) error {
	return b.engine.Shutdown(ctx)
}
