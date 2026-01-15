// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	redigo "github.com/gomodule/redigo/redis"
)

type Redigo struct {
	Client  *redigo.Pool
	options *Options
}

func NewRedigo(options ...Option) *Redigo {
	client := &redigo.Pool{
		MaxIdle:     30,
		MaxActive:   30,
		IdleTimeout: time.Duration(200),
	}
	opts := ApplyOptions(options...)
	if len(opts.RedisHost) == 0 {
		opts.RedisHost = "127.0.0.1:6379"
	}
	client.Dial = func() (redigo.Conn, error) {
		c, err := redigo.Dial("tcp", opts.RedisHost)
		if err != nil {
			return nil, err
		}
		if opts.RedisPassword != "" {
			if _, err := c.Do("AUTH", opts.RedisPassword); err != nil {
				c.Close()
				return nil, err
			}
		}
		return c, err
	}
	client.TestOnBorrow = pingRedis
	return &Redigo{
		Client:  client,
		options: opts,
	}
}

func pingRedis(c redigo.Conn, t time.Time) error {
	_, err := c.Do("PING")
	if err != nil {
		return err
	}
	return nil
}

func (g *Redigo) Get(_ context.Context, key string) (any, error) {
	conn := g.Client.Get()
	defer conn.Close()
	reply, err := redigo.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	return reply, err
}

// GetWithTTL returns data stored from a given key and its corresponding TTL
func (g *Redigo) GetWithTTL(_ context.Context, key string) (any, time.Duration, error) {
	conn := g.Client.Get()
	defer conn.Close()
	reply, err := redigo.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, 0, err
	}
	ttl, err := redigo.Int64(conn.Do("TTL", key))
	if err != nil {
		return reply, 0, err
	}
	if ttl == -2 {
		return nil, 0, errors.Newf("key %s not found", key)
	}
	return reply, time.Duration(ttl), nil
}

// Set defines data in Redis for given key identifier
func (g *Redigo) Set(_ context.Context, key string, value any, options ...Option) error {
	conn := g.Client.Get()
	defer conn.Close()
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	opts := ApplyOptionsWithDefault(g.options, options...)
	_, err = conn.Do("SET", key, data)
	if err != nil {
		return err
	}
	if opts.Expiration > 0 {
		_, err = conn.Do("EXPIRE", key, opts.Expiration.Seconds())
		if err != nil {
			return err
		}
	}
	return nil
}

// Delete removes data from Redis for given key identifier
func (g *Redigo) Delete(_ context.Context, key string) error {
	conn := g.Client.Get()
	defer conn.Close()
	_, err := redigo.Bool(conn.Do("DEL", key))
	return err
}

// Flush resets all data in the store
func (g *Redigo) Flush(_ context.Context) error {
	conn := g.Client.Get()
	defer conn.Close()
	_, err := conn.Do("FLUSHALL")
	return err
}

func (g *Redigo) Ping(_ context.Context) error {
	conn := g.Client.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	return err
}

// Exists checks if a key exists in the store
func (g *Redigo) Exists(key string) bool {
	conn := g.Client.Get()
	defer conn.Close()
	exists, err := redigo.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func (g *Redigo) ReSet(key string, data any, options ...Option) error {
	conn := g.Client.Get()
	defer conn.Close()
	if g.Exists(key) {
		g.Delete(context.Background(), key)
	}
	return g.Set(context.Background(), key, data, options...)
}

func (g *Redigo) LikeDelete(key string) error {
	conn := g.Client.Get()
	defer conn.Close()
	keys, err := redigo.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}
	for _, key := range keys {
		g.Delete(context.Background(), key)
	}
	return nil
}

func (g *Redigo) RPushQueueEnd(queue string, payload any) error {
	conn := g.Client.Get()
	defer conn.Close()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = conn.Do("RPUSH", queue, data)
	return err
}

func (g *Redigo) LPushQueueHeader(queue string, payload any) error {
	conn := g.Client.Get()
	defer conn.Close()
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = conn.Do("LPUSH", queue, data)
	return err
}

// QueueReadHeader read list header
func (g *Redigo) QueueReadHeader(queue string) any {
	conn := g.Client.Get()
	defer conn.Close()
	reply, err := conn.Do("LPOP", queue) // 头部
	if err != nil {
		return nil
	}
	return reply
}

// QueueReadEnd read list end
func (g *Redigo) QueueReadEnd(queue string) any {
	conn := g.Client.Get()
	defer conn.Close()
	reply, err := conn.Do("RPOP", queue) // 尾部
	if err != nil {
		return nil
	}
	return reply
}

func RedigoParse(replay any, t string) (any, error) {
	switch t {
	case "string":
		return redigo.String(replay, nil)
	case "int":
		return redigo.Int(replay, nil)
	case "int64":
		return redigo.Int64(replay, nil)
	case "byte":
		return redigo.Bytes(replay, nil)
	case "bool":
		return redigo.Bool(replay, nil)
	}
	return replay, nil
}
