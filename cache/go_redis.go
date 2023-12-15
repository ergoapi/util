package cache

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	goredis "github.com/redis/go-redis/v9"
)

type GoRedis struct {
	Client  *goredis.Client
	options *Options
}

func NewGoRedis(options ...Option) *GoRedis {
	opts := ApplyOptions(options...)
	if len(opts.RedisHost) == 0 {
		opts.RedisHost = "127.0.0.1:6379"
	}
	client := goredis.NewClient(&goredis.Options{
		Addr:     opts.RedisHost,
		DB:       opts.RedisDB,
		Username: opts.RedisUser,
		Password: opts.RedisPassword,
	})
	return &GoRedis{
		Client:  client,
		options: opts,
	}
}

func (g *GoRedis) Get(key string) (any, error) {
	object, err := g.Client.Get(context.Background(), key).Result()
	if err == goredis.Nil {
		return nil, errors.Newf("key %s not found", key)
	}
	return object, err
}

// GetWithTTL returns data stored from a given key and its corresponding TTL
func (g *GoRedis) GetWithTTL(key string) (any, time.Duration, error) {
	object, err := g.Client.Get(context.Background(), key).Result()
	if err == goredis.Nil {
		return nil, 0, errors.Newf("key %s not found", key)
	}
	if err != nil {
		return nil, 0, err
	}

	ttl, err := g.Client.TTL(context.Background(), key).Result()
	if err != nil {
		return nil, 0, err
	}

	return object, ttl, err
}

// Set defines data in Redis for given key identifier
func (g *GoRedis) Set(key string, value any, options ...Option) error {
	opts := ApplyOptionsWithDefault(g.options, options...)

	err := g.Client.Set(context.Background(), key, value, opts.Expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Delete removes data from Redis for given key identifier
func (g *GoRedis) Delete(key string) error {
	_, err := g.Client.Del(context.Background(), key).Result()
	return err
}

// Clear resets all data in the store
func (g *GoRedis) Flush() error {
	return g.Client.FlushAll(context.Background()).Err()
}

func (g *GoRedis) Ping() error {
	return g.Client.Ping(context.Background()).Err()
}
