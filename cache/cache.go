// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package cache provides a unified caching interface with multiple backend support.
package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	GetWithTTL(ctx context.Context, key string) (any, time.Duration, error)
	Set(ctx context.Context, key string, value any, options ...Option) error
	Delete(ctx context.Context, key string) error
	Flush(ctx context.Context) error
	Ping(ctx context.Context) error
}

var (
	Instance Cache
)

func Get(ctx context.Context, key string) (any, error) {
	return Instance.Get(ctx, key)
}

func GetWithTTL(ctx context.Context, key string) (any, time.Duration, error) {
	return Instance.GetWithTTL(ctx, key)
}

func Set(ctx context.Context, key string, value any, options ...Option) error {
	return Instance.Set(ctx, key, value, options...)
}

func Delete(ctx context.Context, key string) error {
	return Instance.Delete(ctx, key)
}

func Flush(ctx context.Context) error {
	return Instance.Flush(ctx)
}

func Ping(ctx context.Context) error {
	return Instance.Ping(ctx)
}
