// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type user struct {
	ts     int64
	tokens uint
}

type inMemoryStoreType struct {
	rate   int64
	limit  uint
	data   *sync.Map
	mu     sync.Map // per-key mutex
	skip   func(ctx *gin.Context) bool
	cancel context.CancelFunc
}

func (s *inMemoryStoreType) getMutex(key string) *sync.Mutex {
	mu, _ := s.mu.LoadOrStore(key, &sync.Mutex{})
	return mu.(*sync.Mutex)
}

func (s *inMemoryStoreType) Limit(key string, c *gin.Context) Info {
	mu := s.getMutex(key)
	mu.Lock()
	defer mu.Unlock()

	var u user
	m, ok := s.data.Load(key)
	if !ok {
		u = user{time.Now().Unix(), s.limit}
	} else {
		u = m.(user)
	}
	if u.ts+s.rate <= time.Now().Unix() {
		u.tokens = s.limit
	}
	if s.skip != nil && s.skip(c) {
		return Info{
			Limit:         s.limit,
			RateLimited:   false,
			ResetTime:     time.Now().Add(time.Duration((s.rate - (time.Now().Unix() - u.ts)) * time.Second.Nanoseconds())),
			RemainingHits: u.tokens,
		}
	}
	if u.tokens <= 0 {
		return Info{
			Limit:         s.limit,
			RateLimited:   true,
			ResetTime:     time.Now().Add(time.Duration((s.rate - (time.Now().Unix() - u.ts)) * time.Second.Nanoseconds())),
			RemainingHits: 0,
		}
	}
	u.tokens--
	u.ts = time.Now().Unix()
	s.data.Store(key, u)
	return Info{
		Limit:         s.limit,
		RateLimited:   false,
		ResetTime:     time.Now().Add(time.Duration((s.rate - (time.Now().Unix() - u.ts)) * time.Second.Nanoseconds())),
		RemainingHits: u.tokens,
	}
}

// Close stops the background cleanup goroutine.
func (s *inMemoryStoreType) Close() {
	if s.cancel != nil {
		s.cancel()
	}
}

type InMemoryOptions struct {
	// the user can make Limit amount of requests every Rate
	Rate time.Duration
	// the amount of requests that can be made every Rate
	Limit uint
	// a function that returns true if the request should not count toward the rate limit
	Skip func(*gin.Context) bool
}

// InMemoryStore creates a new in-memory rate limit store.
// Call Close() on the returned store when done to stop the background cleanup goroutine.
func InMemoryStore(options *InMemoryOptions) *inMemoryStoreType {
	ctx, cancel := context.WithCancel(context.Background())
	data := &sync.Map{}
	store := &inMemoryStoreType{
		rate:   int64(options.Rate.Seconds()),
		limit:  options.Limit,
		data:   data,
		skip:   options.Skip,
		cancel: cancel,
	}
	go clearInBackground(ctx, data, store.rate)
	return store
}

func clearInBackground(ctx context.Context, data *sync.Map, rate int64) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			data.Range(func(k, v any) bool {
				if v.(user).ts+rate <= time.Now().Unix() {
					data.Delete(k)
				}
				return true
			})
		}
	}
}
