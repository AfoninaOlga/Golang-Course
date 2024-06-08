package handler

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConcurrencyLimiter_Add(t *testing.T) {
	limiter := NewConcurrencyLimiter(2)
	limiter.Add()
	limiter.Add()
	assert.Equal(t, 2, len(limiter.sem))
}

func TestConcurrencyLimiter_Remove(t *testing.T) {
	limiter := NewConcurrencyLimiter(2)
	limiter.Add()
	assert.Equal(t, 1, len(limiter.sem))
	limiter.Remove()
	assert.Equal(t, 0, len(limiter.sem))
}

func TestRateLimiter_Allow(t *testing.T) {
	limiter := NewRateLimiter(context.Background(), 1, time.Second)
	allowed := limiter.Allow("user")
	assert.True(t, allowed)
	allowed = limiter.Allow("user")
	assert.False(t, allowed)
}

func TestRateLimiter_deleteExpired(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	limiter := NewRateLimiter(ctx, 1, time.Second/8)
	limiter.Allow("user")
	time.Sleep(time.Second / 10)
	assert.Contains(t, limiter.clientMap, "user")
	time.Sleep(time.Second / 10)
	assert.NotContains(t, limiter.clientMap, "user")
	cancel()
}

func TestRateLimiter_Allow2(t *testing.T) {

}
