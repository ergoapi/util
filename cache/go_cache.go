// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package cache

import (
	"time"

	"github.com/cockroachdb/errors"
	gocache "github.com/patrickmn/go-cache"
)

type GoCache struct {
	Client  *gocache.Cache
	options *Options
}

func NewGoCache(options ...Option) *GoCache {
	opts := ApplyOptions(options...)
	return &GoCache{
		Client:  gocache.New(opts.Expiration, opts.CleanupInterval),
		options: opts,
	}
}

func (g *GoCache) Get(key string) (any, error) {
	var err error
	value, found := g.Client.Get(key)
	if !found {
		err = errors.Newf("key %s not found", key)
	}
	return value, err
}

func (g *GoCache) GetWithTTL(key string) (any, time.Duration, error) {
	var err error
	value, t, found := g.Client.GetWithExpiration(key)
	if !found {
		err = errors.Newf("key %s not found", key)
		return value, 0, err
	}
	return value, time.Until(t), nil
}

func (g *GoCache) Set(key string, value any, options ...Option) error {
	opts := ApplyOptions(options...)
	if opts == nil {
		opts = g.options
	}
	g.Client.Set(key, value, opts.Expiration)
	return nil
}

func (g *GoCache) Delete(key string) error {
	g.Client.Delete(key)
	return nil
}

func (g *GoCache) Flush() error {
	g.Client.Flush()
	return nil
}

func (g *GoCache) Ping() error {
	return nil
}
