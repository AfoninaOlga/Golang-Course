package handler

import (
	"context"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	"github.com/AfoninaOlga/xkcd/internal/core/port/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuth_ServiceError(t *testing.T) {
	as := new(mocks.AuthService)
	as.On("GetUserByToken", mock.Anything, mock.Anything).Return(nil, resErr)
	req, err := http.NewRequest("POST", "/update", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	handler := Auth(true, as, nil)
	handler(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuth_NotAdmin(t *testing.T) {
	as := new(mocks.AuthService)
	as.On("GetUserByToken", mock.Anything, mock.Anything).Return(&domain.User{IsAdmin: false}, nil)
	req, err := http.NewRequest("POST", "/update", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	handler := Auth(true, as, nil)
	handler(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestAuth(t *testing.T) {
	as := new(mocks.AuthService)
	as.On("GetUserByToken", mock.Anything, mock.Anything).Return(&domain.User{IsAdmin: false, Name: "user"}, nil)
	req, err := http.NewRequest("GET", "", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	handler := Auth(false, as, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
	}))
	handler(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestConcurrencyLimiting(t *testing.T) {
	limiter := NewConcurrencyLimiter(10)
	req, err := http.NewRequest("GET", "", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	handler := ConcurrencyLimiting(limiter, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	handler(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestRateLimiting_NoUser(t *testing.T) {
	limiter := NewRateLimiter(context.Background(), 10, time.Minute)
	handler := RateLimiting(limiter, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	req, err := http.NewRequest("GET", "", nil)
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	handler(recorder, req)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
}

func TestRateLimiting_IPUser(t *testing.T) {
	limiter := NewRateLimiter(context.Background(), 10, time.Minute)
	handler := RateLimiting(limiter, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	req, err := http.NewRequest("GET", "", nil)
	req.RemoteAddr = "127.0.0.1:4242"
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	handler(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestRateLimiting_ManyRequests(t *testing.T) {
	limiter := NewRateLimiter(context.Background(), 1, time.Minute)
	handler := RateLimiting(limiter, http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {}))
	req, err := http.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:4242"
	assert.NoError(t, err)
	recorder1 := httptest.NewRecorder()
	recorder2 := httptest.NewRecorder()
	handler(recorder1, req)
	handler(recorder2, req)
	assert.Contains(t, []int{recorder1.Code, recorder2.Code}, http.StatusTooManyRequests)
}

func TestRateLimiting_ContextUser(t *testing.T) {
	limiter := NewRateLimiter(context.Background(), 1, time.Minute)
	handler := RateLimiting(limiter, http.HandlerFunc(func(_ http.ResponseWriter, req *http.Request) {}))
	req, err := http.NewRequest("GET", "/", nil)
	ctx := context.WithValue(context.Background(), ctxKey("user"), domain.User{Name: "user"})
	assert.NoError(t, err)
	recorder := httptest.NewRecorder()
	handler(recorder, req.WithContext(ctx))
	assert.Equal(t, http.StatusOK, recorder.Code)
}
