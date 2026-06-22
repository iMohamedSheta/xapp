// Package xqueue is a generic in-process worker-pool engine. It knows
// nothing about notifications, listeners, or any domain type — it only
// moves values of type T from Enqueue to a Handler, with N workers,
// graceful shutdown, and optional delayed dispatch via a timer goroutine.
package goroutine_backend

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"
)

const (
	DefaultWorkers   = 10
	DefaultQueueSize = 512
)

// ErrStopped is returned by Enqueue once the engine has begun shutting
// down. It is also the error a blocked Enqueue call sees if Shutdown is
// triggered while it's waiting for queue space.
var ErrStopped = errors.New("xqueue: engine is stopped")

type Options struct {
	Workers   int
	QueueSize int
	Logger    Logger
}

func (o Options) workers() int {
	if o.Workers <= 0 {
		return DefaultWorkers
	}
	return o.Workers
}

func (o Options) queueSize() int {
	if o.QueueSize <= 0 {
		return DefaultQueueSize
	}
	return o.QueueSize
}

func (o Options) logger() Logger {
	if o.Logger == nil {
		return slog.Default()
	}
	return o.Logger
}

type item[T any] struct {
	ctx context.Context
	val T
}

// Handler processes one dequeued item. Engine never logs or interprets
// the returned error — it's fire-and-forget by design. The wrapping
// backend's handler closure is the one that knows the right fields to
// log on failure, so it owns that decision.
type Handler[T any] func(ctx context.Context, val T) error

type Engine[T any] struct {
	queue chan item[T]
	stop  chan struct{} // closed once on Shutdown; lets blocked sends/timers bail early
	once  sync.Once

	mu     sync.RWMutex // guards closed; see Enqueue/Shutdown for the gating proof
	closed bool

	pending sync.WaitGroup // in-flight Enqueue calls (immediate sends + delayed timers)
	wg      sync.WaitGroup // worker goroutines only, fixed size, added once in New

	handler Handler[T]
	logger  Logger
}

func New[T any](opts Options, handler Handler[T]) *Engine[T] {
	e := &Engine[T]{
		queue:   make(chan item[T], opts.queueSize()),
		stop:    make(chan struct{}),
		handler: handler,
		logger:  opts.logger(),
	}
	for range opts.workers() {
		e.wg.Add(1)
		go e.worker()
	}
	return e
}

// Enqueue hands val to a worker immediately, or — if scheduleAt is in the
// future — defers it to a per-item timer goroutine.
//
// The closed check and pending.Add are done together under mu.RLock so
// they can never race against Shutdown's mu.Lock: once Shutdown observes
// no readers and sets closed = true, no later Enqueue can sneak an Add in
// after pending.Wait() has already been (or is being) evaluated.
func (e *Engine[T]) Enqueue(ctx context.Context, val T, scheduleAt *time.Time) error {
	e.mu.RLock()
	if e.closed {
		e.mu.RUnlock()
		e.logger.Warn("xqueue: engine is stopped")
		return ErrStopped
	}
	e.pending.Add(1)
	e.mu.RUnlock()

	it := item[T]{ctx: ctx, val: val}

	if scheduleAt == nil || !scheduleAt.After(time.Now()) {
		defer e.pending.Done()
		select {
		case e.queue <- it:
			return nil
		case <-e.stop:
			e.logger.Warn("xqueue: engine is stopped")
			return ErrStopped
		}
	}

	go func() {
		defer e.pending.Done()
		t := time.NewTimer(time.Until(*scheduleAt))
		defer t.Stop()
		select {
		case <-t.C:
			// Safe unguarded send: e.queue can't be closed while this
			// goroutine's pending.Add(1) is still outstanding.
			e.queue <- it
		case <-e.stop:
			e.logger.Warn("xqueue: engine is stopped; canceling scheduled item")
		}
	}()
	return nil
}

// Shutdown stops accepting new work, lets in-flight Enqueue calls and
// pending timers resolve, then closes the queue so workers drain whatever
// was already buffered and exit. It returns early with ctx.Err() if ctx
// is done before that finishes.
func (e *Engine[T]) Shutdown(ctx context.Context) error {
	e.once.Do(func() {
		e.mu.Lock()
		e.closed = true
		e.mu.Unlock()

		close(e.stop) // wakes any timer/send currently blocked waiting on it

		go func() {
			e.pending.Wait() // no more sends can start once closed; wait for the rest to finish
			close(e.queue)   // safe now — workers range to completion and exit
		}()
	})

	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (e *Engine[T]) worker() {
	defer e.wg.Done()
	for it := range e.queue {
		if err := e.handler(it.ctx, it.val); err != nil {
			e.logger.Error("xqueue: handler error", "error", err)
		}
	}
}
