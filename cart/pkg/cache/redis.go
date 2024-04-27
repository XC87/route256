package cache

import (
	"context"
	"errors"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"route256.ozon.ru/project/cart/internal/config"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Redis struct {
	cache *cache.Cache
	ttl   time.Duration
}

func NewRedis(ctx context.Context, config *config.Config, wg *sync.WaitGroup) *Redis {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{"shard0": config.RedisHost},
	})

	redisCache := cache.New(&cache.Options{
		Redis:        ring,
		StatsEnabled: true,
	})

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		err := ring.Close()
		if err != nil {
			zap.L().Error("cant close redis connections", zap.Error(err))
		}
	}()

	return &Redis{cache: redisCache, ttl: config.RedisTTL}
}

func (r Redis) StartMonitorHitMiss(ctx context.Context, registerer prometheus.Registerer) {
	lastHits := r.cache.Stats().Hits
	lastMiss := r.cache.Stats().Misses
	cacheHits := promauto.With(registerer).NewCounter(prometheus.CounterOpts{
		Name: "cache_hits",
		Help: "The total number of cache  hits",
	})
	cacheMisses := promauto.With(registerer).NewCounter(prometheus.CounterOpts{
		Name: "cache_misses",
		Help: "The total number of cache misses",
	})

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				currentHits := r.cache.Stats().Hits
				currentMiss := r.cache.Stats().Misses

				cacheHits.Add(float64(currentHits - lastHits))
				cacheMisses.Add(float64(currentMiss - lastMiss))

				lastHits = currentHits
				lastMiss = currentMiss
			}
		}
	}()
}

func (r Redis) Get(ctx context.Context, key string) (value string, err error) {
	err = r.cache.Get(ctx, key, &value)
	if errors.Is(err, cache.ErrCacheMiss) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return value, nil
}

func (r Redis) Set(key string, value string) error {
	return r.cache.Set(&cache.Item{
		Key:   key,
		Value: value,
		TTL:   r.ttl,
	})
}

func (r Redis) Invalidate(ctx context.Context, key string) error {
	return r.cache.Delete(ctx, key)
}
