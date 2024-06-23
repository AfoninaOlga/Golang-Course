package handler

import (
	"context"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type Client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	limit      int
	clientMap  map[string]*Client
	expireTime time.Duration
	mtx        sync.Mutex
}

type ConcurrencyLimiter struct {
	sem chan struct{}
}

func NewConcurrencyLimiter(limit int) *ConcurrencyLimiter {
	return &ConcurrencyLimiter{
		sem: make(chan struct{}, limit),
	}
}

func (cl *ConcurrencyLimiter) Add() {
	cl.sem <- struct{}{}
}

func (cl *ConcurrencyLimiter) Remove() {
	<-cl.sem
}

func NewRateLimiter(ctx context.Context, limit int, expireTime time.Duration) *RateLimiter {
	rl := &RateLimiter{limit: limit, clientMap: make(map[string]*Client), expireTime: expireTime}
	ticker := time.NewTicker(expireTime)
	go func() {
		for {
			select {
			case <-ticker.C:
				rl.deleteExpired()
			case <-ctx.Done():
				return
			}
		}
	}()
	return rl
}

func (rl *RateLimiter) Allow(client string) bool {
	rl.mtx.Lock()
	defer rl.mtx.Unlock()

	if c, ok := rl.clientMap[client]; ok {
		c.lastSeen = time.Now()
	} else {
		rl.clientMap[client] = &Client{limiter: rate.NewLimiter(rate.Limit(rl.limit), rl.limit), lastSeen: time.Now()}
	}
	return rl.clientMap[client].limiter.Allow()
}

func (rl *RateLimiter) deleteExpired() {
	rl.mtx.Lock()
	defer rl.mtx.Unlock()
	for key, client := range rl.clientMap {
		if time.Since(client.lastSeen) > rl.expireTime {
			delete(rl.clientMap, key)
		}
	}
}
