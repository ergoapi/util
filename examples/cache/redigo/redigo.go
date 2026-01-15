// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package main

import (
	"context"
	"log"
	"time"

	"github.com/ergoapi/util/cache"
)

func main() {
	ctx := context.Background()
	redigo := cache.NewRedigo(
		cache.WithRedisHost("10.143.5.228:6379"),
		cache.WithRedisPassword("oxX5eh5OQuivaigheexahge5Nahth3xe"),
	)
	if err := redigo.Ping(ctx); err != nil {
		return
	}
	log.Println("cache is ready")
	if err := redigo.Set(ctx, "key", "value"); err != nil {
		panic(err)
	}
	value, err := redigo.Get(ctx, "key")
	if err != nil {
		panic(err)
	}
	value, _ = cache.RedigoParse(value, "string")
	log.Println("value: ", value)
	if err := redigo.Set(ctx, "key1", "value1", cache.WithExpiration(time.Second*15)); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 10)
	for {
		if value, t, err := redigo.GetWithTTL(ctx, "key1"); err != nil {
			log.Println("key1 not found")
			break
		} else {
			value, _ = cache.RedigoParse(value, "string")
			log.Println("key1 value:", value, t)
			if value, ttl, err := redigo.GetWithTTL(ctx, "key"); err != nil {
				log.Println("key not found")
			} else {
				value, _ = cache.RedigoParse(value, "string")
				log.Println("key value:", value, " ttl: ", ttl.Seconds())
			}
			time.Sleep(time.Second * 2)
		}
	}

	// if err := redigo.Flush(ctx); err != nil {
	// 	panic(err)
	// }
}
