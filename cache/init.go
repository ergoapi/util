// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package cache

func InitGoCache(options ...Option) {
	Instance = NewGoCache(options...)
}

func InitGoRedis(options ...Option) {
	Instance = NewGoRedis(options...)
}

func InitGoRedisCluster(options ...Option) {
	Instance = NewGoRedisCluster(options...)
}

func InitRedigo(options ...Option) {
	Instance = NewRedigo(options...)
}
