// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package cache provides a unified caching interface with multiple backend support.
package cache

import "time"

type Cache interface {
	Get(key string) (any, error)
	GetWithTTL(key string) (any, time.Duration, error)
	Set(key string, value any, options ...Option) error
	Delete(key string) error
	Flush() error
	Ping() error
}

var (
	Instance Cache
)

func Get(key string) (any, error) {
	return Instance.Get(key)
}

func GetWithTTL(key string) (any, time.Duration, error) {
	return Instance.GetWithTTL(key)
}

func Set(key string, value any, options ...Option) error {
	return Instance.Set(key, value, options...)
}

func Delete(key string) error {
	return Instance.Delete(key)
}

func Flush() error {
	return Instance.Flush()
}

func Ping() error {
	return Instance.Ping()
}
