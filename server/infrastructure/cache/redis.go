package cache

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

var once sync.Once

type RedisImpl struct {
	RedisClientRing *redis.Ring
	CacheClient     *cache.Cache
}

func NewRedisClient() *RedisImpl {
	log.Println("Initialize Redis connection")
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"127.0.0.1": ":6380",
		},
		Password: "eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81",
	})

	cache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return &RedisImpl{
		RedisClientRing: ring,
		CacheClient:     cache,
	}
}

func (c RedisImpl) Close() error {
	return c.RedisClientRing.Close()
}
