package redis_backend

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/imohamedsheta/xapp/pkg/bus"
	"github.com/redis/go-redis/v9"
)

// Logger is a minimal structured logger interface.
// Compatible with slog.Logger, zap.SugaredLogger, and similar.
type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Warn(msg string, fields ...any)
}

// QueueHandler mirrors the goroutine backend's dispatch contract —
// satisfied by *bus.Bus via HandleTask.
type QueueHandler interface {
	HandleTask(ctx context.Context, taskKey string, payload []byte) error
}

const (
	fieldTaskKey = "task_key"
	fieldPayload = "payload"
)

type Options struct {
	StreamKey    string // required, e.g. bus.Prefix() with the trailing ':' trimmed
	Group        string // required, consumer group name
	Consumer     string // required, this process's consumer name within the group
	Concurrency  int
	BlockTimeout time.Duration
	ClaimMinIdle time.Duration

	// MaxDeliveries caps how many times a single entry gets reclaimed and
	// retried before reclaimOnce gives up and routes it to DeadLetterStream
	// instead. 0 means "use the default" (5) — there is no "unlimited"
	// option; an uncapped retry loop on a poison message is exactly the gap
	// this field exists to close.
	MaxDeliveries int64

	// DeadLetterStream is where entries land after exceeding MaxDeliveries.
	// Defaults to StreamKey + ":dead" if empty. Nothing consumes this
	// stream automatically — it's a holding pen for manual inspection.
	DeadLetterStream string

	// IdempotencyTTL, if > 0, records a "handled" marker per stream entry
	// ID after a successful Handle, keyed before the ack. A redelivery of
	// that same entry (e.g. the ack itself was lost) then skips re-running
	// Handle and just retries the ack. Set well above
	// ClaimMinIdle * (MaxDeliveries + 2) so the marker can't expire while
	// retries are still in flight. 0 disables this check entirely.
	IdempotencyTTL time.Duration

	Logger Logger
}

func (o Options) logger() Logger {
	if o.Logger == nil {
		return slog.Default()
	}
	return o.Logger
}

type RedisStreamsBackend struct {
	client           *redis.Client
	streamKey        string
	group            string
	consumer         string
	conc             int
	block            time.Duration
	minIdle          time.Duration
	maxDeliveries    int64
	deadLetterStream string
	idempotencyTTL   time.Duration
	logger           Logger

	mu      sync.RWMutex
	handler QueueHandler

	stop chan struct{}
	wg   sync.WaitGroup
}

func New(client *redis.Client, opts Options) (*RedisStreamsBackend, error) {
	if client == nil {
		return nil, errors.New("bus - redis streams: nil client")
	}
	if opts.StreamKey == "" {
		return nil, errors.New("bus - redis streams: StreamKey is required")
	}
	if opts.Group == "" {
		return nil, errors.New("bus - redis streams: Group is required")
	}
	if opts.Consumer == "" {
		return nil, errors.New("bus - redis streams: Consumer is required")
	}

	deadLetterStream := opts.DeadLetterStream
	if deadLetterStream == "" {
		deadLetterStream = opts.StreamKey + ":dead"
	}

	b := &RedisStreamsBackend{
		client:           client,
		streamKey:        opts.StreamKey,
		group:            opts.Group,
		consumer:         opts.Consumer,
		conc:             orDefaultInt(opts.Concurrency, 1),
		block:            orDefaultDur(opts.BlockTimeout, 5*time.Second),
		minIdle:          orDefaultDur(opts.ClaimMinIdle, 30*time.Second),
		maxDeliveries:    orDefaultInt64(opts.MaxDeliveries, 5),
		deadLetterStream: deadLetterStream,
		idempotencyTTL:   opts.IdempotencyTTL,
		logger:           opts.logger(),
		stop:             make(chan struct{}),
	}

	if err := b.ensureGroup(context.Background()); err != nil {
		return nil, err
	}
	return b, nil
}

func orDefaultInt(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}

func orDefaultInt64(v, def int64) int64 {
	if v <= 0 {
		return def
	}
	return v
}

func orDefaultDur(v, def time.Duration) time.Duration {
	if v <= 0 {
		return def
	}
	return v
}

// ensureGroup creates the consumer group + stream (MKSTREAM) if it doesn't
// exist yet. Safe to call on every boot — BUSYGROUP just means it's already
// there.
func (b *RedisStreamsBackend) ensureGroup(ctx context.Context) error {
	err := b.client.XGroupCreateMkStream(ctx, b.streamKey, b.group, "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("bus - redis streams: create group: %w", err)
	}
	return nil
}

// SetHandler wires the dispatch side. Call before Start.
func (b *RedisStreamsBackend) SetHandler(h QueueHandler) {
	b.mu.Lock()
	b.handler = h
	b.mu.Unlock()
}

// Enqueue implements bus.QueueBackend — XADD onto the shared stream, tagged
// with taskKey as a field, so every consumer reads one stream and routes by
// key (the same model AsynqHandler uses for asynq's prefix-based mux).
func (b *RedisStreamsBackend) Enqueue(ctx context.Context, _ bus.Listener, taskKey string, payload []byte) error {
	err := b.client.XAdd(ctx, &redis.XAddArgs{
		Stream: b.streamKey,
		Values: map[string]any{
			fieldTaskKey: taskKey,
			fieldPayload: payload,
		},
	}).Err()
	if err != nil {
		return fmt.Errorf("bus - redis streams: xadd: %w", err)
	}
	return nil
}

// Start launches Concurrency consumer loops plus one reclaim loop. Call
// once, after SetHandler. Returns immediately; loops run until ctx is
// cancelled or Shutdown is called.
func (b *RedisStreamsBackend) Start(ctx context.Context) {
	for i := 0; i < b.conc; i++ {
		b.wg.Add(1)
		go b.consume(ctx)
	}
	b.wg.Add(1)
	go b.reclaimLoop(ctx)
}

func (b *RedisStreamsBackend) consume(ctx context.Context) {
	defer b.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-b.stop:
			return
		default:
		}

		res, err := b.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    b.group,
			Consumer: b.consumer,
			Streams:  []string{b.streamKey, ">"},
			Count:    10,
			Block:    b.block,
		}).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) || errors.Is(err, context.Canceled) {
				continue // BLOCK timed out, nothing new — loop again
			}
			b.logger.Error("bus - redis streams: xreadgroup", "error", err)
			time.Sleep(time.Second)
			continue
		}

		for _, stream := range res {
			for _, msg := range stream.Messages {
				b.dispatch(ctx, msg)
			}
		}
	}
}

func (b *RedisStreamsBackend) dispatch(ctx context.Context, msg redis.XMessage) {
	taskKey, _ := msg.Values[fieldTaskKey].(string)
	payload := payloadBytes(msg.Values[fieldPayload])
	dedupeKey := "bus:done:" + msg.ID

	if b.idempotencyTTL > 0 {
		already, err := b.client.Exists(ctx, dedupeKey).Result()
		if err != nil {
			// Fail open: a Redis blip on the check itself shouldn't block
			// real work. Worst case here is one avoidable duplicate run,
			// which idempotent listeners absorb anyway.
			b.logger.Error("bus - redis streams: dedupe check failed", "id", msg.ID, "error", err)
		} else if already > 0 {
			// Handle already ran and succeeded for this exact entry; we're
			// only here because that ack got lost. Don't re-run Handle.
			b.logger.Warn("bus - redis streams: skipping already-handled redelivery", "id", msg.ID, "task_key", taskKey)
			b.ackWithRetry(ctx, msg.ID)
			return
		}
	}

	b.mu.RLock()
	h := b.handler
	b.mu.RUnlock()

	if h == nil {
		b.logger.Error("bus - redis streams: no handler registered", "task_key", taskKey, "id", msg.ID)
		return // leave unacked — reclaimOnce retries it once a handler exists
	}

	if err := h.HandleTask(ctx, taskKey, payload); err != nil {
		b.logger.Error("bus - redis streams: handler error", "task_key", taskKey, "id", msg.ID, "error", err)
		return // leave unacked — reclaimOnce retries it, up to MaxDeliveries
	}

	if b.idempotencyTTL > 0 {
		// Write the marker before acking, not after: if the process dies
		// between these two lines, the next redelivery sees the marker,
		// skips Handle, and just retries the ack — which is exactly the
		// outcome we want. Doing it in the other order would let a crash
		// after a successful ack-but-failed-marker-write through with no
		// protection at all.
		if err := b.client.Set(ctx, dedupeKey, 1, b.idempotencyTTL).Err(); err != nil {
			b.logger.Error("bus - redis streams: dedupe write failed", "id", msg.ID, "error", err)
		}
	}

	b.ackWithRetry(ctx, msg.ID)
}

// ackWithRetry retries XAck itself a few times before giving up. Handle
// already succeeded by the time this is called — losing the ack to a
// transient network error and giving up immediately just turns a clean
// success into an avoidable duplicate redelivery later.
func (b *RedisStreamsBackend) ackWithRetry(ctx context.Context, id string) {
	const maxAttempts = 3
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err = b.client.XAck(ctx, b.streamKey, b.group, id).Err(); err == nil {
			return
		}
		if attempt < maxAttempts {
			time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
		}
	}
	// Still failed after retries. reclaimOnce will eventually redeliver
	// this — the idempotency marker (if enabled) is what keeps that
	// redelivery from re-running Handle.
	b.logger.Error("bus - redis streams: xack failed after retries", "id", id, "error", err)
}

// payloadBytes handles both the []byte we sent and the string go-redis may
// hand back, depending on driver version / RESP encoding.
func payloadBytes(v any) []byte {
	switch t := v.(type) {
	case []byte:
		return t
	case string:
		return []byte(t)
	default:
		return nil
	}
}

// reclaimLoop periodically claims entries that some consumer read but never
// acked (crashed, panicked, or died mid-handle) and re-dispatches them on
// this consumer. Without this, a crashed consumer's in-flight messages sit
// in the pending list forever.
func (b *RedisStreamsBackend) reclaimLoop(ctx context.Context) {
	defer b.wg.Done()
	ticker := time.NewTicker(b.minIdle)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-b.stop:
			return
		case <-ticker.C:
			b.reclaimOnce(ctx)
		}
	}
}

// reclaimOnce inspects entries idle past minIdle, dead-letters anything
// that's exceeded MaxDeliveries, and claims + redispatches the rest.
//
// Pagination note: this processes at most one page (Count entries) per
// tick rather than looping to exhaustion. Idle time only grows for entries
// left behind, so anything not caught this tick is still caught on a later
// tick — trading a bit of redelivery latency under a large backlog for not
// having to get Redis Stream ID range pagination exactly right.
func (b *RedisStreamsBackend) reclaimOnce(ctx context.Context) {
	const pageSize = 100

	pending, err := b.client.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: b.streamKey,
		Group:  b.group,
		Idle:   b.minIdle,
		Start:  "-",
		End:    "+",
		Count:  pageSize,
	}).Result()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			b.logger.Error("bus - redis streams: xpendingext", "error", err)
		}
		return
	}
	if len(pending) == 0 {
		return
	}

	retryIDs := make([]string, 0, len(pending))
	for _, p := range pending {
		if p.RetryCount > b.maxDeliveries {
			b.deadLetter(ctx, p.ID, p.RetryCount)
			continue
		}
		retryIDs = append(retryIDs, p.ID)
	}
	if len(retryIDs) == 0 {
		return
	}

	msgs, err := b.client.XClaim(ctx, &redis.XClaimArgs{
		Stream:   b.streamKey,
		Group:    b.group,
		Consumer: b.consumer,
		MinIdle:  b.minIdle,
		Messages: retryIDs,
	}).Result()
	if err != nil {
		b.logger.Error("bus - redis streams: xclaim", "error", err)
		return
	}
	for _, msg := range msgs {
		b.dispatch(ctx, msg)
	}
}

// deadLetter moves a poison entry off the live stream's pending list onto
// DeadLetterStream instead of retrying it forever. XPendingExt only reports
// metadata (ID, consumer, idle, retry count), not the original fields, so
// the payload is looked up via XRange before being copied over.
func (b *RedisStreamsBackend) deadLetter(ctx context.Context, id string, retryCount int64) {
	entries, err := b.client.XRange(ctx, b.streamKey, id, id).Result()
	if err != nil || len(entries) == 0 {
		b.logger.Error("bus - redis streams: dead-letter lookup failed", "id", id, "error", err)
		// Ack anyway: without the fields we can't retry it meaningfully
		// either, and leaving it pending forever is strictly worse than
		// dropping it with a loud log line.
		_ = b.client.XAck(ctx, b.streamKey, b.group, id).Err()
		return
	}
	orig := entries[0]

	err = b.client.XAdd(ctx, &redis.XAddArgs{
		Stream: b.deadLetterStream,
		Values: map[string]any{
			fieldTaskKey:  orig.Values[fieldTaskKey],
			fieldPayload:  orig.Values[fieldPayload],
			"retry_count": retryCount,
			"original_id": id,
		},
	}).Err()
	if err != nil {
		b.logger.Error("bus - redis streams: dead-letter xadd failed", "id", id, "error", err)
		return // leave it pending — better to retry once more than ack-and-lose it
	}

	if err := b.client.XAck(ctx, b.streamKey, b.group, id).Err(); err != nil {
		b.logger.Error("bus - redis streams: dead-letter xack failed", "id", id, "error", err)
		return
	}
	b.logger.Warn("bus - redis streams: dead-lettered after max retries",
		"id", id, "retry_count", retryCount, "task_key", orig.Values[fieldTaskKey])
}

// Shutdown signals all loops to stop and waits for them to exit, or returns
// ctx.Err() if ctx is cancelled first.
func (b *RedisStreamsBackend) Shutdown(ctx context.Context) error {
	close(b.stop)
	done := make(chan struct{})
	go func() {
		b.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
