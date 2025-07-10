package cache

import (
	"context"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/simple-auth-roles/internal/config"
)

// CacheService defines cache operations for temporary data storage
type CacheService interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type cacheService struct {
	client redis.UniversalClient
	logger *slog.Logger
}

func NewCacheService(cfg *config.Config, logger *slog.Logger) CacheService {
	var client redis.UniversalClient

	if cfg.Redis.URL != "" {
		// Use Redis if configured
		opts, err := redis.ParseURL(cfg.Redis.URL)
		if err != nil {
			logger.Warn("Invalid Redis URL, using in-memory cache", "error", err)
			return &memoryCacheService{
				data:   make(map[string]cacheItem),
				logger: logger.With("cache", "memory"),
			}
		}
		client = redis.NewClient(opts)

		// Test connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			logger.Warn("Redis connection failed, using in-memory cache", "error", err)
			return &memoryCacheService{
				data:   make(map[string]cacheItem),
				logger: logger.With("cache", "memory"),
			}
		}

		logger.Info("Connected to Redis cache")
		return &cacheService{
			client: client,
			logger: logger.With("cache", "redis"),
		}
	}

	// Use in-memory cache as fallback
	logger.Info("Using in-memory cache (Redis not configured)")
	return &memoryCacheService{
		data:   make(map[string]cacheItem),
		logger: logger.With("cache", "memory"),
	}
}

func (c *cacheService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *cacheService) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *cacheService) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
