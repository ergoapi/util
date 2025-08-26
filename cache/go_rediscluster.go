// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package cache

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	goredis "github.com/redis/go-redis/v9"
)

type GoRedisCluster struct {
	Client  *goredis.ClusterClient
	options *Options
}

func NewGoRedisCluster(options ...Option) *GoRedisCluster {
	opts := ApplyOptions(options...)
	if len(opts.RedisHost) == 0 {
		opts.RedisHost = "127.0.0.1:6379,"
	}
	client := goredis.NewClusterClient(&goredis.ClusterOptions{
		Addrs:    strings.Split(opts.RedisHost, ","),
		Username: opts.RedisUser,
		Password: opts.RedisPassword,
	})
	return &GoRedisCluster{
		Client:  client,
		options: opts,
	}
}

func (g *GoRedisCluster) Get(key string) (any, error) {
	object, err := g.Client.Get(context.Background(), key).Result()
	if err == goredis.Nil {
		return nil, errors.Newf("key %s not found", key)
	}
	return object, err
}

// GetWithTTL returns data stored from a given key and its corresponding TTL
func (g *GoRedisCluster) GetWithTTL(key string) (any, time.Duration, error) {
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
func (g *GoRedisCluster) Set(key string, value any, options ...Option) error {
	opts := ApplyOptionsWithDefault(g.options, options...)

	err := g.Client.Set(context.Background(), key, value, opts.Expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Delete removes data from Redis for given key identifier
func (g *GoRedisCluster) Delete(key string) error {
	_, err := g.Client.Del(context.Background(), key).Result()
	return err
}

// Flush resets all data in the store
func (g *GoRedisCluster) Flush() error {
	return g.Client.FlushAll(context.Background()).Err()
}

func (g *GoRedisCluster) Ping() error {
	return g.Client.Ping(context.Background()).Err()
}
