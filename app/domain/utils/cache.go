package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/iMohamedSheta/xerr"
	"github.com/imohamedsheta/xapp/app/domain/enums"
	"github.com/imohamedsheta/xapp/app/x"
	"github.com/redis/go-redis/v9"
)

const (
	RememberForever time.Duration = 0
)

// CacheRemember tries to get the value from cache. If not found, it runs the fallback,
// stores the result in Redis as JSON, and returns it.
func CacheRemember[T any](ctx context.Context, key string, ttl time.Duration, fallback func() (T, error)) (T, error) {
	// Get Redis connection internally
	cache := x.Cache().Redis

	var value T
	var err error

	// Try to get from cache
	cachedData, err := cache.Get(ctx, key).Result()
	if err == nil {
		// Found in cache, unmarshal it
		if err := json.Unmarshal([]byte(cachedData), &value); err == nil {
			return value, nil
		}
		// If unmarshal fails, we'll proceed to fallback
	} else if err != redis.Nil {
		// Redis error (not just missing key)
		return value, fmt.Errorf("failed to fetch from Redis: %w", err)
	}

	// Not found or unmarshal failed, run fallback
	value, err = fallback()
	if err != nil {
		return value, err
	}

	// Marshal and cache the value
	jsonData, err := json.Marshal(value)
	if err != nil {
		return value, fmt.Errorf("failed to marshal value for caching: %w", err)
	}

	if err := cache.Set(ctx, key, string(jsonData), ttl).Err(); err != nil {
		fmt.Printf("Warning: Failed to cache value for key %s: %v\n", key, err)
	}

	return value, nil
}

// CacheForget removes a key from the cache.
func CacheForget(ctx context.Context, key string) error {
	cache := x.Cache().Redis
	return cache.Del(ctx, key).Err()
}

func GetSessionSettingsFromCache(c *gin.Context) (map[string]string, error) {
	key, err := GetAuthSessionSettingsCacheKey(c)
	if err != nil {
		return nil, err
	}

	if r := x.Cache(); r != nil {
		return r.Redis.HGetAll(c, key).Result()
	}

	return nil, xerr.New("Cache client Redis not configured", enums.XErrServerError, nil)
}
