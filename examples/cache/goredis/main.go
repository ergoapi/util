package main

import (
	"log"
	"time"

	"github.com/ergoapi/util/cache"
)

func main() {
	cache.InitGoRedis(cache.WithRedisDB(0),
		cache.WithRedisHost("localhost:6379"),
		cache.WithRedisPassword("password123"),
	)
	if err := cache.Ping(); err != nil {
		return
	}
	log.Println("cache is ready")
	if err := cache.Set("key", "value"); err != nil {
		panic(err)
	}
	value, err := cache.Get("key")
	if err != nil {
		panic(err)
	}
	log.Println("value:", value)
	if err := cache.Set("key1", "value1", cache.WithExpiration(time.Second*15)); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 10)
	for {
		if value, t, err := cache.GetWithTTL("key1"); err != nil {
			log.Println("key1 not found")
			break
		} else {
			log.Println("key1 value:", value, t)
			if value, err := cache.Get("key"); err != nil {
				log.Println("key not found")
			} else {
				log.Println("key value:", value)
			}
			time.Sleep(time.Second * 2)
		}
	}

	if err := cache.Flush(); err != nil {
		panic(err)
	}
}
