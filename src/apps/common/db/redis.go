package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	redisPoolSize     = 1000             // 可有可无
	redisMinIdleConns = 10               // 可有可无
	redisIdleTimeout  = 30 * time.Second // 可有可无
)

type ICacheUtil interface {
	GetCmdTimeout() time.Duration
	GetCacheDB() *redis.Client

	Get(key string) (string, error)
	Set(key string, value string, expiration time.Duration) error
}

type _CacheUtil struct {
	cacheDB         *redis.Client
	cacheCmdTimeout time.Duration
}

func NewCacheUtil(redisAddr, redisPassword string, redisDB int, redisCmdTimeoutSec time.Duration) (ICacheUtil, error) {
	if redisCmdTimeoutSec == 0 {
		panic("NewCacheUtil must have redisCmdTimeoutSec ")
	}
	opt := redis.Options{
		Addr:         redisAddr,
		Password:     redisPassword,
		DB:           redisDB,
		PoolSize:     redisPoolSize,
		MinIdleConns: redisMinIdleConns,
		IdleTimeout:  redisIdleTimeout,
	}
	r := redis.NewClient(&opt)
	if err := r.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &_CacheUtil{
		cacheDB:         r,
		cacheCmdTimeout: redisCmdTimeoutSec,
	}, nil
}

func (u *_CacheUtil) GetCmdTimeout() time.Duration {
	return u.cacheCmdTimeout
}

func (u *_CacheUtil) GetCacheDB() *redis.Client {
	return u.cacheDB
}

func (u *_CacheUtil) Get(key string) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), u.GetCmdTimeout())
	return u.GetCacheDB().Get(ctx, key).Result()
}

func (u *_CacheUtil) Set(key string, value string, expiration time.Duration) error {
	ctx, _ := context.WithTimeout(context.Background(), u.GetCmdTimeout())
	return u.GetCacheDB().Set(ctx, key, value, expiration).Err()
}
