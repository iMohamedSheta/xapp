package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xapp/bootstrap/support"
	"github.com/imohamedsheta/xioc"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// WebsocketHooks holds app-specific dependencies injected from the app layer.
type WebsocketHooks struct {
	// RegisterRoutes registers app WebSocket routes onto the given router group.
	RegisterRoutes func(group *gin.RouterGroup, auth gin.HandlerFunc)
	// BuildAuthMiddleware constructs the auth middleware from config.
	BuildAuthMiddleware func() gin.HandlerFunc
	// Middleware holds the standard gin middleware factories.
	Middleware GinMiddleware
}

// GinMiddleware holds gin middleware factories so bootstrap doesn't import app/middleware.
type GinMiddleware struct {
	Recovery func() gin.HandlerFunc
	Logger   func() gin.HandlerFunc
	CORS     func() gin.HandlerFunc
}

// Shutdown gracefully closes all loaded services.
func Shutdown() {
	support.PrintWarning("Shutting down services...")
	xioc.AppContainer().ShutdownAll()
	support.PrintSuccess("All services shut down properly")
}

// RunHttp starts the HTTP/HTTPS server with graceful shutdown.
// handler is the root http.Handler (e.g. routes.RegisterRoutes()).
func RunHttp(handler http.Handler) {
	cfg := x.Config()
	shutdownTimeout := cfg.GetDuration("app.shutdown_timeout", 10*time.Second)
	bindAddress := cfg.GetString("app.bind_address", "0.0.0.0")
	bindPort := cfg.GetString("app.bind_port", "8080")
	url := cfg.GetString("app.url", "http://localhost")

	httpsCfg := cfg.GetMap("app.https", map[string]any{"enabled": false})
	httpsEnabled, _ := httpsCfg["enabled"].(bool)
	certFile, _ := httpsCfg["certFile"].(string)
	keyFile, _ := httpsCfg["keyFile"].(string)
	httpsPort, _ := httpsCfg["port"].(string)

	srv := &http.Server{
		Addr:    bindAddress + ":" + bindPort,
		Handler: handler,
	}

	if httpsEnabled {
		srv.Addr = bindAddress + ":" + httpsPort
		localCertPath := x.Storage().DefaultDisk().Path(certFile)
		localKeyPath := x.Storage().DefaultDisk().Path(keyFile)
		support.SafeGo(func() {
			support.PrintInfo("Starting HTTPS server on https://" + url + ":" + httpsPort + "/health ...")
			if err := srv.ListenAndServeTLS(localCertPath, localKeyPath); err != nil && err != http.ErrServerClosed {
				x.Logger().Error("HTTPS Server error", zap.Error(err))
				support.PrintErr(err)
				os.Exit(1)
			}
		})
	} else {
		support.SafeGo(func() {
			support.PrintInfo("Starting HTTP server on http://" + url + ":" + bindPort + "/health ...")
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				x.Logger().Error("HTTP Server error", zap.Error(err))
				os.Exit(1)
			}
		})
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	support.PrintWarning("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		x.Logger().Error("Server forced to shutdown", zap.Error(err))
		os.Exit(1)
	}

	Shutdown()
	support.PrintSuccess("Server exited properly")
}

// RunWorker starts the asynq task worker.
// handlers is the full map of task type → handler (already built by the caller, including orchestrator).
// websocketQueues controls whether to use the websocket queue config key instead of the default.
func RunWorker(handlers map[string]asynq.Handler, websocketQueues bool) {
	cfg := x.Config()

	if !cfg.GetBool("queue.enabled", true) {
		log.Fatal("Queue is disabled, cannot start consumer")
	}

	redisCfg := cfg.GetMap("redis.connections.queue", nil)
	if redisCfg == nil {
		log.Fatal("redis connection config is missing for queue consumer")
	}

	redisOpt, err := convertRedisConfigToOptions(redisCfg)
	if err != nil {
		log.Fatal("Failed to get Redis options:", err)
	}

	concurrency := cfg.GetInt("queue.consumer.concurrency", 10)

	var queuesRaw any
	if websocketQueues {
		queuesRaw, _ = cfg.Get("queue.consumer.websocket_queues")
	} else {
		queuesRaw, _ = cfg.Get("queue.consumer.queues")
	}
	queues := convertQueuePriorities(queuesRaw)

	serverConfig := asynq.Config{
		Concurrency:    concurrency,
		Queues:         queues,
		RetryDelayFunc: asynq.DefaultRetryDelayFunc,
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			if cfg.GetBool("queue.consumer.logging.log_failed_tasks", true) {
				retryCount, _ := asynq.GetRetryCount(ctx)
				maxRetry, _ := asynq.GetMaxRetry(ctx)
				taskID, _ := asynq.GetTaskID(ctx)
				x.Logger("queue_log").Error("Task processing failed",
					zap.String("task_type", task.Type()),
					zap.String("task_id", taskID),
					zap.Int("retry_count", retryCount),
					zap.Int("max_retry", maxRetry),
					zap.Error(err),
				)
			}
		}),
	}

	server := asynq.NewServer(asynq.RedisClientOpt{
		Addr:      redisOpt.Addr,
		Username:  redisOpt.Username,
		Password:  redisOpt.Password,
		DB:        redisOpt.DB,
		PoolSize:  redisOpt.PoolSize,
		TLSConfig: redisOpt.TLSConfig,
	}, serverConfig)

	mux := asynq.NewServeMux()
	log.Printf("Registered %d task handlers", len(handlers))
	for taskType, handler := range handlers {
		log.Printf("Registering task handler: %s", taskType)
		mux.Handle(taskType, wrapHandlerWithLogging(handler, taskType))
	}

	log.Printf("Starting task worker with %d concurrency...", concurrency)
	log.Printf("Queue priorities: %+v", queues)

	if err := server.Run(mux); err != nil {
		log.Fatal("Task server error:", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Printf("Shutting down task server...")
	server.Shutdown()
	Shutdown()
}

// RunWebsocket starts the dedicated WebSocket HTTP server.
// It also starts a websocket-specific asynq worker in the background.
// hooks provides all app-specific route registration and middleware.
func RunWebsocket(handlers map[string]asynq.Handler, hooks WebsocketHooks) {
	// Start the websocket asynq worker in background
	support.SafeGo(func() {
		RunWorker(handlers, true)
	})

	cfg := x.Config()
	shutdownTimeout := cfg.GetDuration("app.shutdown_timeout", 10*time.Second)
	bindAddress := cfg.GetString("websocket.bind_address", "0.0.0.0")
	bindPort := cfg.GetString("websocket.bind_port", "8081")

	r := gin.Default()
	r.Use(hooks.Middleware.Recovery())
	r.Use(hooks.Middleware.Logger())
	r.Use(hooks.Middleware.CORS())

	authMiddleware := hooks.BuildAuthMiddleware()
	wsGroup := r.Group("/")
	hooks.RegisterRoutes(wsGroup, authMiddleware)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "websocket"})
	})

	srv := &http.Server{
		Addr:    bindAddress + ":" + bindPort,
		Handler: r,
	}

	support.SafeGo(func() {
		support.PrintInfo("Starting WebSocket server on ws://" + bindAddress + ":" + bindPort + " ...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			x.Logger().Error("WebSocket server error", zap.Error(err))
			os.Exit(1)
		}
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	support.PrintWarning("Shutting down WebSocket server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		x.Logger().Error("WebSocket server forced to shutdown", zap.Error(err))
		os.Exit(1)
	}

	Shutdown()
	support.PrintSuccess("WebSocket server exited properly")
}

// RunScheduler starts the periodic task scheduler.
// registerSchedule is the app-level function that registers all scheduled tasks.
func RunScheduler(registerSchedule func()) {
	schedulerClient := x.Scheduler()
	if schedulerClient == nil {
		log.Println("Scheduler not found in container")
		return
	}

	registerSchedule()

	support.PrintInfo("Starting periodic task scheduler...")
	if err := schedulerClient.Scheduler.Run(); err != nil {
		log.Printf("Scheduler error: %v", err)
	}
}

// wrapHandlerWithLogging decorates a handler with success/failure logging.
func wrapHandlerWithLogging(handler asynq.Handler, taskType string) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, task *asynq.Task) error {
		start := time.Now()
		err := handler.ProcessTask(ctx, task)
		duration := time.Since(start)

		if err != nil {
			return err
		}

		if x.Config().GetBool("queue.consumer.logging.log_success", false) {
			taskID, _ := asynq.GetTaskID(ctx)
			x.Logger("queue_log").Info("Task completed successfully",
				zap.String("task_type", taskType),
				zap.String("task_id", taskID),
				zap.Duration("duration", duration),
			)
		}

		return nil
	})
}

func convertQueuePriorities(raw any) map[string]int {
	result := make(map[string]int)
	queuesMap, ok := raw.(map[string]any)
	if !ok {
		return result
	}
	for k, v := range queuesMap {
		switch val := v.(type) {
		case int:
			result[k] = val
		case int64:
			result[k] = int(val)
		case float64:
			result[k] = int(val)
		case string:
			if parsed, err := strconv.Atoi(val); err == nil {
				result[k] = parsed
			}
		}
	}
	return result
}

func convertRedisConfigToOptions(cfg map[string]any) (*redis.Options, error) {
	if isActive, ok := cfg["active"].(bool); ok && !isActive {
		return nil, errors.New("redis connection is not active")
	}

	if rawURL, ok := cfg["url"].(string); ok && rawURL != "" {
		opt, err := redis.ParseURL(rawURL)
		if err != nil {
			return nil, fmt.Errorf("invalid redis url: %w", err)
		}
		return opt, nil
	}

	host := cfg["host"].(string)
	port := cfg["port"].(int)
	password, _ := cfg["password"].(string)
	db, _ := cfg["database"].(int)
	poolSize, _ := cfg["pool_size"].(int)

	opt := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	}
	if poolSize > 0 {
		opt.PoolSize = poolSize
	}
	if timeoutStr, ok := cfg["timeout"].(string); ok && timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			opt.DialTimeout = timeout
		}
	}
	return opt, nil
}
