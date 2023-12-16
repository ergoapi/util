package main

import (
	"log"
	"time"

	"github.com/ergoapi/util/cache"
)

func main() {
	redigo := cache.NewRedigo(
		cache.WithRedisHost("10.143.5.228:6379"),
		cache.WithRedisPassword("oxX5eh5OQuivaigheexahge5Nahth3xe"),
	)
	if err := redigo.Ping(); err != nil {
		return
	}
	log.Println("cache is ready")
	if err := redigo.Set("key", "value"); err != nil {
		panic(err)
	}
	value, err := redigo.Get("key")
	if err != nil {
		panic(err)
	}
	value, _ = cache.RedigoParse(value, "string")
	log.Println("value: ", value)
	if err := redigo.Set("key1", "value1", cache.WithExpiration(time.Second*15)); err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 10)
	for {
		if value, t, err := redigo.GetWithTTL("key1"); err != nil {
			log.Println("key1 not found")
			break
		} else {
			value, _ = cache.RedigoParse(value, "string")
			log.Println("key1 value:", value, t)
			if value, ttl, err := redigo.GetWithTTL("key"); err != nil {
				log.Println("key not found")
			} else {
				value, _ = cache.RedigoParse(value, "string")
				log.Println("key value:", value, " ttl: ", ttl.Seconds())
			}
			time.Sleep(time.Second * 2)
		}
	}

	// if err := redigo.Flush(); err != nil {
	// 	panic(err)
	// }
}
