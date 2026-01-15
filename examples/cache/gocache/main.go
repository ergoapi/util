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
	cache.InitGoCache()
	if err := cache.Ping(ctx); err != nil {
		return
	}
	log.Println("cache is ready")
	if err := cache.Set(ctx, "key", "value"); err != nil {
		panic(err)
	}
	value, err := cache.Get(ctx, "key")
	if err != nil {
		panic(err)
	}
	log.Println("value:", value)
	if err := cache.Set(ctx, "key1", "value1", cache.WithExpiration(time.Second*15)); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 10)
	for {
		if value, t, err := cache.GetWithTTL(ctx, "key1"); err != nil {
			log.Println("key1 not found")
			break
		} else {
			log.Println("key1 value:", value, t)
			if value, err := cache.Get(ctx, "key"); err != nil {
				log.Println("key not found")
			} else {
				log.Println("key value:", value)
			}
			time.Sleep(time.Second * 2)
		}
	}

	if err := cache.Flush(ctx); err != nil {
		panic(err)
	}
}
