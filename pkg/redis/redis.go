package redis

import (
	"context"
	"fmt"
	"time"

	goRedis "github.com/go-redis/redis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// Nil this for standard redis nil
	Nil = "redis: nil"
)

// Adapter this interface for redis
type Adapter interface {
	Ping(ctx context.Context) (string, error)
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	Subscribe(ctx context.Context, channels ...string) *goRedis.PubSub
	Publish(ctx context.Context, channel string, data string) (int64, error)
}

// Config config for redis
type Config struct {
	Addrs []string
	DB    int
}

// NewRedis this for new redis
func NewRedis(serviceName string, cfg Config) Adapter {
	redisClient := new(Redis)
	cacheOptions := &goRedis.UniversalOptions{
		Addrs:    cfg.Addrs,
		PoolSize: 10,
	}
	redisClient.serviceName = fmt.Sprintf("%s.redis", serviceName)
	redisClient.config = cfg
	redisClient.client = goRedis.NewUniversalClient(cacheOptions)
	if _, err := redisClient.Ping(context.Background()); err != nil {
		log.Fatal(fmt.Sprintf("error connect redis err: %s", err.Error()))
	}
	return redisClient
}

// Redis this for struct redis
type Redis struct {
	serviceName string
	client      goRedis.UniversalClient
	config      Config
}

// Get this for get redis
func (r *Redis) Get(ctx context.Context, key string) (string, error) {

	stringCMD := r.client.Get(key)
	if stringCMD.Err() == goRedis.Nil {
		return "", errors.New(Nil)
	} else if stringCMD.Err() != nil {
		return "", stringCMD.Err()
	}
	return stringCMD.String(), nil
}

// Set this for set redis
func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (string, error) {

	return r.client.Set(key, value, expiration).Result()
}

// Exists this for check exist redis
func (r *Redis) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Exists(keys...).Result()
}

// Ping check connection redis
func (r *Redis) Ping(ctx context.Context) (string, error) {
	return r.client.Ping().Result()
}

// Incr this for increment redis
func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {

	return r.client.Incr(key).Result()
}

// Decr this for decrement redis
func (r *Redis) Decr(ctx context.Context, key string) (int64, error) {
	return r.client.Decr(key).Result()
}

// Del this for delete keys
func (r *Redis) Del(ctx context.Context, keys ...string) (int64, error) {
	return r.client.Del(keys...).Result()
}

// Subscribe subsribe message
func (r *Redis) Subscribe(ctx context.Context, channels ...string) *goRedis.PubSub {
	return r.client.Subscribe(channels...)
}

// Publish  message
func (r *Redis) Publish(ctx context.Context, channel string, data string) (int64, error) {
	return r.client.Publish(channel, data).Result()
}
