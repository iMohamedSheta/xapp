package load

import (
	"context"
	"errors"
	"os"

	"github.com/hibiken/asynq"
	"github.com/imohamedsheta/xapp/app/domain/utils"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/imohamedsheta/xioc"
	"github.com/redis/go-redis/v9"
)

func InitRedisCache(c *xioc.Container) {
	// Load cache connection
	err := xioc.Singleton(c, func(c *xioc.Container) (*x.CacheClient, error) {
		defaultConnection := x.Config().GetString("redis.default", "default")
		cacheClient, err := cacheClient("redis.connections." + defaultConnection)
		if err != nil {
			return nil, err
		}
		return x.NewCacheClient(cacheClient), nil
	})

	if err != nil {
		errMsg := "Failed to load redis cache module in the ioc container : " + err.Error()
		x.Logger().Error(errMsg)
		utils.PrintErr(errMsg)
		os.Exit(1)
	}
}

func InitRedisQueue(c *xioc.Container) {
	// 1. Get Redis Options
	ops, err := utils.GetRedisQueueClientOptionsForAsynq("redis.connections.queue")
	if err != nil {
		logRedisError("Failed to get redis queue options", err)
	}

	// 2. Load queue client connection
	err = xioc.Singleton(c, func(c *xioc.Container) (*x.QueueClient, error) {
		client := asynq.NewClient(ops)
		if err := client.Ping(); err != nil {
			return nil, err
		}
		return x.NewQueueClient(client), nil
	})
	if err != nil {
		logRedisError("Failed to load redis queue client", err)
	}

	// 3. Load scheduler connection
	err = xioc.Singleton(c, func(c *xioc.Container) (*x.SchedulerClient, error) {
		scheduler := asynq.NewScheduler(ops, nil)
		return x.NewSchedulerClient(scheduler), nil
	})
	if err != nil {
		logRedisError("Failed to load redis scheduler client", err)
	}
}

func logRedisError(msg string, err error) {
	errMsg := msg + " : " + err.Error()
	x.Logger().Error(errMsg)
	utils.PrintErr(errMsg)
	os.Exit(1)
}

// loads the redis cache connection
func cacheClient(key string) (*redis.Client, error) {
	cfg := x.Config().GetMap(key, nil)
	if cfg == nil {
		return nil, errors.New("redis connection config not found for " + key)
	}

	options, err := utils.ConvertRedisConfigToOptions(cfg)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		errMsg := "Failed to load redis cache module in the ioc container : " + err.Error()
		x.Logger().Error(errMsg)
		utils.PrintErr(errMsg)
		os.Exit(1)
	}

	return client, nil
}

// loads the queue connection from the config
func queueClient(key string) (*asynq.Client, error) {
	ops, err := utils.GetRedisQueueClientOptionsForAsynq(key)
	if err != nil {
		return nil, err
	}

	client := asynq.NewClient(ops)

	if err := client.Ping(); err != nil {
		errMsg := "Failed to load redis queue module in the ioc container : " + err.Error()
		x.Logger().Error(errMsg)
		utils.PrintErr(errMsg)
		os.Exit(1)
	}

	return client, nil
}
