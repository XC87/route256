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
	latestHits := r.cache.Stats().Hits
	latestMisses := r.cache.Stats().Misses

	go func() {
		cacheHits := promauto.With(registerer).NewCounter(prometheus.CounterOpts{
			Name: "cart_list_cache_hits",
			Help: "The total number of cart list hits",
		})
		cacheMisses := promauto.With(registerer).NewCounter(prometheus.CounterOpts{
			Name: "cart_list_cache_misses",
			Help: "The total number of cart list misses",
		})

		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				currentHits := r.cache.Stats().Hits
				currentMisses := r.cache.Stats().Misses

				cacheHits.Add(float64(currentHits - latestHits))
				cacheMisses.Add(float64(currentMisses - latestMisses))

				latestHits = currentHits
				latestMisses = currentMisses
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
